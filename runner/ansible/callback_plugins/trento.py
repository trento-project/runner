"""
Trento Callback executor.

:author: xarbulu
:organization: SUSE Linux GmbH
:contact: xarbulu@suse.com

:since: 2021-09-16
"""

import logging
import os
import yaml
import requests

from ansible.plugins.callback import CallbackBase

TRENTO_TEST_LABEL_KEY = "trento_labels"
TRENTO_TEST_LABEL = "test"
TEST_RESULT_TASK_NAME = "set_test_result"
TEST_INCLUDE_TASK_NAME = "run_checks"
CHECK_ID = "id"

EXECUTION_COMPLETED_EVENT = "execution_completed"


class ExecutionResults(object):
    """
    Object to store and user the execution results

    Result example:

    {
      "cluster_id": "cluster1",
      "hosts": [
        {
          "host_id": "host1",
          "name": "host1",
          "reachable": true,
          "msg": "some message",
          "results": [
            {
              "check_id: "check1",
              "result": "passing",
              "msg": "some message"
            },
            {
              "check_id: "check1",
              "result": "warning",
              "msg": "some message"
            }
          ]
        }
      ]
    }
    """
    def __init__(self):
        self._logger = logging.getLogger(__name__)
        self.cluster = None

    def initialize_cluster(self, cluster_id):
        """
        Initialize the cluster object
        """
        self.cluster = Cluster(cluster_id)

    def add_host(self, host_id, state, msg=""):
        """
        Add the host state. Reachable or Unreachable
        """
        self.cluster.add_host(host_id, state, msg)

    def add_result(self, host_id, check_id, result, msg=""):
        """
        Add check result
        """
        self.cluster.add_result(host_id, check_id, result, msg)

    def to_dict(self):
        """
        Transform to dictionary
        """
        return self.cluster.to_dict()


class Cluster(object):
    """
    Cluster data object
    """

    def __init__(self, cluster_id):
        self.cluster_id = cluster_id
        self.hosts = []

    def add_host(self, host_id, state, msg=""):
        """
        Set the host state. Reachable or Unreachable
        """
        for host in self.hosts:
            if host.host_id == host_id:
                host.reachable = state
                host.msg = msg
                break
        else:
            self.hosts.append(Host(host_id, state, msg))

    def add_result(self, host_id, check_id, result, msg=""):
        """
        Add check result
        """
        for host in self.hosts:
            if host.host_id == host_id:
                host.add_result(check_id, result, msg)
                break

    def to_dict(self):
        """
        Transform to dictionary
        """
        return {
            "cluster_id": self.cluster_id,
            "hosts": [host.to_dict() for host in self.hosts]
        }


class Host(object):
    """
    Host data object
    """

    def __init__(self, host_id, reachable, msg):
        self.host_id = host_id
        self.results = []
        self.reachable = reachable
        self.msg = msg

    def add_result(self, check_id, result, msg=""):
        """
        Add check result
        """
        # Check if a result already exists
        # Due how ansible callbacks system works, we might get same check results twice
        # where the 2nd result is false, as it gives the results of `set_test_result` task
        # when the check has failed due abnormal behaviours
        for result_item in self.results:
            if result_item.check_id == check_id:
                break
        else:
            self.results.append(CheckResult(check_id, result, msg))

    def to_dict(self):
        """
        Transform to dictionary
        """
        return {
            "host_id": self.host_id,
            "reachable": self.reachable,
            "results": [result.to_dict() for result in self.results],
            "msg": self.msg
        }


class CheckResult(object):
    """
    Check result data object
    """

    def __init__(self, check_id, result, msg):
        self.check_id = check_id
        self.result = result
        self.msg = msg

    def to_dict(self):
        """
        Transform to dictionary
        """
        return {
            "check_id": str(self.check_id),
            "result": self.result,
            "msg": self.msg
        }


class CallbackModule(CallbackBase):
    """
    Trento Callback module
    """
    CALLBACK_VERSION = 2.0
    CALLBACK_TYPE = 'aggregate'
    CALLBACK_NAME = 'trento'

    def __init__(self):
        super(CallbackModule, self).__init__()
        self.play = None
        self.execution_results = ExecutionResults()
        self._callbacks_url = os.getenv('TRENTO_CALLBACKS_URL')
        self._execution_id = os.getenv('TRENTO_EXECUTION_ID')

    def v2_playbook_on_start(self, _):
        """
        On start callback
        """
        self._display.banner("Trento callbacks plugin loaded")

    def v2_playbook_on_play_start(self, play):
        """
        On Play start callback
        """
        self.play = play
        play_vars = self._all_vars()
        for _, host_data in play_vars["hostvars"].items():
            for group in host_data["group_names"]:
                self.execution_results.initialize_cluster(group)
                break

    def v2_runner_on_ok(self, result):
        """
        On task Ok
        """
        if self._is_check_include_loop(result):
            self._store_skipped(result)
            return

        if not self._is_test_result(result):
            return

        host = result._host.get_name()
        task_vars = self._all_vars(host=result._host, task=result._task)

        test_result = result._task_fields["args"]["test_result"]
        self.execution_results.add_host(host, True)
        self.execution_results.add_result(host, task_vars[CHECK_ID], test_result)

    def v2_runner_on_failed(self, result, ignore_errors):
        """
        On task Failed
        """
        host = result._host.get_name()
        task_vars = self._all_vars(host=result._host, task=result._task)

        if CHECK_ID not in task_vars:
            return

        msg = result._check_key("msg")
        self.execution_results.add_host(host, True)
        self.execution_results.add_result(host, task_vars[CHECK_ID], "critical", msg)

    def v2_runner_on_skipped(self, result):
        """
        On task Skipped
        """
        if self._is_check_include_loop(result):
            self._store_skipped(result)

    def v2_runner_on_unreachable(self, result):
        """
        On task Unreachable
        """
        host = result._host.get_name()
        msg = result._check_key("msg")
        self.execution_results.add_host(host, False, msg)

    def v2_playbook_on_stats(self, _stats):
        """
        Post results at the end of the execution
        """
        if not self._is_test_execution():
            return

        self._display.banner("Publishing Trento results")
        self._post_results()

    def _all_vars(self, host=None, task=None):
        """
        Get task vars

        host and task need to be specified in case 'magic variables' (host vars, group vars, etc)
        need to be loaded as well
        """
        return self.play.get_variable_manager().get_vars(
            play=self.play,
            host=host,
            task=task
        )

    def _is_test_execution(self):
        """
        Check if the current execution is a trento test execution
        """
        play_vars = self._all_vars()
        if TRENTO_TEST_LABEL_KEY not in play_vars or \
                 TRENTO_TEST_LABEL not in play_vars[TRENTO_TEST_LABEL_KEY]:
            self._display.banner("Not running a Trento test execution")
            return False
        return True

    def _is_test_result(self, result):
        """
        Check if the current task is a test result
        """
        if (result._task_fields.get("action") == "set_fact") and \
                (result._task_fields.get("name") == TEST_RESULT_TASK_NAME):
            return True
        return False

    def _is_check_include_loop(self, result):
        """
        Check if the current task is the checks include loop task
        """
        if (result._task_fields.get("action") == "include_role") and \
                (result._task_fields.get("name") == TEST_INCLUDE_TASK_NAME):
            return True
        return False

    def _store_skipped(self, result):
        """
        Store skipped checks
        """
        host = result._host.get_name()

        for check_result in result._result["results"]:
            skipped = check_result.get("skipped", False)
            if skipped:
                with open(os.path.join(
                    check_result["check_item"]["path"], "defaults/main.yml")) as file_ptr:

                    data = yaml.load(file_ptr, Loader=yaml.Loader)
                    check_id = data[CHECK_ID]

                self.execution_results.add_host(host, True)
                self.execution_results.add_result(host, check_id, "skipped")

    def _post_results(self):
        """
        Post results to the trento web api server
        """
        callback_event = {
            "execution_id": self._execution_id,
            "event": EXECUTION_COMPLETED_EVENT,
            "payload": self.execution_results.to_dict()
        }
        response = requests.post(self._callbacks_url, json=callback_event)
        self._display.banner(
            "Results of execution {} published. Return code is: {}".format(
                self._execution_id, response.status_code))

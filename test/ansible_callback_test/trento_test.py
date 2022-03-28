"""
Unitary tests for the ansible callbacks trento.py
"""

# pylint:disable=C0103,C0111,W0212,W0611

import os
import sys
import unittest

sys.path.insert(
    0, os.path.abspath(
        os.path.join(
            os.path.dirname(__file__),
            '../../runner/ansible/callback_plugins')))


import trento

class TestTrentoCallbacks(unittest.TestCase):
    """
    Unitary tests for trento.py.
    """


    def test_initialize_cluster(self):
        result = trento.ExecutionResults()
        result.initialize_cluster("cluster1")
        result.initialize_cluster("cluster2")
        result.initialize_cluster("cluster1")

        expected_result = [
            {
                "cluster_id": "cluster1",
                "hosts": []
            },
            {
                "cluster_id": "cluster2",
                "hosts": []
            }
        ]
        assert expected_result == result.to_dict()

    def test_add_result(self):
        result = trento.ExecutionResults()
        result.initialize_cluster("cluster1")
        result.initialize_cluster("cluster2")

        result.add_result("cluster1", "host1", "other check", "passing", "check message")

        result.add_host("cluster1", "host1", True, "some message")
        result.add_host("cluster1", "host2", True, "other message")
        result.add_host("cluster1", "host3", False, "unreachable")
        result.add_host("cluster2", "host3", True)
        result.add_host("cluster3", "host1", True)

        result.add_result("cluster1", "host1", "check1", "passing", "check message")
        result.add_result("cluster1", "host1", "check2", "critical", "critical message")
        result.add_result("cluster1", "host2", "check1", "passing", "check message")
        result.add_result("cluster2", "host3", "check1", "passing")
        result.add_result("cluster2", "host3", "check3", "warning")
        result.add_result("cluster3", "host1", "check1", "passing")

        expected_result = [
            {
                "cluster_id": "cluster1",
                "hosts": [
                    {
                        "host_id": "host1",
                        "reachable": True,
                        "msg": "some message",
                        "results": [
                            {
                                "check_id": "check1",
                                "result": "passing",
                                "msg": "check message"
                            },
                            {
                                "check_id": "check2",
                                "result": "critical",
                                "msg": "critical message"
                            }
                        ]
                    },
                    {
                        "host_id": "host2",
                        "reachable": True,
                        "msg": "other message",
                        "results": [
                            {
                                "check_id": "check1",
                                "result": "passing",
                                "msg": "check message"
                            }
                        ]
                    },
                    {
                        "host_id": "host3",
                        "reachable": False,
                        "msg": "unreachable",
                        "results": []
                    }
                ]
            },
            {
                "cluster_id": "cluster2",
                "hosts": [
                    {
                        "host_id": "host3",
                        "reachable": True,
                        "msg": "",
                        "results": [
                            {
                                "check_id": "check1",
                                "result": "passing",
                                "msg": ""
                            },
                            {
                                "check_id": "check3",
                                "result": "warning",
                                "msg": ""
                            }
                        ]
                    }
                ]
            }
        ]
        assert expected_result == result.to_dict()

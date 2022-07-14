# Changelog

## [1.0.1](https://github.com/trento-project/runner/tree/1.0.1) (2022-07-14)

[Full Changelog](https://github.com/trento-project/runner/compare/1.0.0...1.0.1)

### Added

- Include suse container packaging [\#57](https://github.com/trento-project/runner/pull/57) (@arbulu89)

### Fixed

- added \s\* whitespace between "token" and ":" [\#78](https://github.com/trento-project/runner/pull/78) (@ksanjeet)
- Fix inconsisten vendoring issue due the self replace [\#73](https://github.com/trento-project/runner/pull/73) (@arbulu89)
- Fix corosync file regexp to be more robust [\#65](https://github.com/trento-project/runner/pull/65) (@arbulu89)
- Add obs disk constraints to avoid out of space issue [\#63](https://github.com/trento-project/runner/pull/63) (@arbulu89)

### Other Changes

- Bump github.com/swaggo/swag from 1.8.2 to 1.8.3 [\#75](https://github.com/trento-project/runner/pull/75) (@dependabot[bot])
- Bump actions/setup-python from 3 to 4 [\#71](https://github.com/trento-project/runner/pull/71) (@dependabot[bot])
- Bump github.com/gin-gonic/gin from 1.7.7 to 1.8.1 [\#70](https://github.com/trento-project/runner/pull/70) (@dependabot[bot])
- Bump github.com/stretchr/testify from 1.7.1 to 1.7.2 [\#69](https://github.com/trento-project/runner/pull/69) (@dependabot[bot])
- Bump github.com/spf13/viper from 1.11.0 to 1.12.0 [\#67](https://github.com/trento-project/runner/pull/67) (@dependabot[bot])
- Bump github.com/vektra/mockery/v2 from 2.12.1 to 2.12.3 [\#66](https://github.com/trento-project/runner/pull/66) (@dependabot[bot])
- Bump github.com/swaggo/swag from 1.8.1 to 1.8.2 [\#64](https://github.com/trento-project/runner/pull/64) (@dependabot[bot])
- Bump docker/login-action from 1.14.1 to 2 [\#61](https://github.com/trento-project/runner/pull/61) (@dependabot[bot])
- Bump docker/metadata-action from 3.8.0 to 4.0.1 [\#60](https://github.com/trento-project/runner/pull/60) (@dependabot[bot])
- Bump docker/setup-buildx-action from 1 to 2 [\#59](https://github.com/trento-project/runner/pull/59) (@dependabot[bot])
- Bump docker/build-push-action from 2 to 3 [\#58](https://github.com/trento-project/runner/pull/58) (@dependabot[bot])
- Restore binaries upload [\#56](https://github.com/trento-project/runner/pull/56) (@dottorblaster)

## [1.0.0](https://github.com/trento-project/runner/tree/1.0.0) (2022-04-29)

[Full Changelog](https://github.com/trento-project/runner/compare/7da8894a9f0f423aaa9fcdd3ff20f07788e9b13c...1.0.0)

### Added

- Update execution requested and completed event payload [\#28](https://github.com/trento-project/runner/pull/28) (@arbulu89)
- Update get catalog endpoint to return 204 when not ready [\#14](https://github.com/trento-project/runner/pull/14) (@arbulu89)
- Flatten checks catalog ouptut [\#13](https://github.com/trento-project/runner/pull/13) (@arbulu89)
- Update execution payload format [\#11](https://github.com/trento-project/runner/pull/11) (@arbulu89)
- Execute playbook with the incoming clusters data [\#9](https://github.com/trento-project/runner/pull/9) (@arbulu89)
- Add webhook callbacks usage [\#8](https://github.com/trento-project/runner/pull/8) (@arbulu89)
- Add execution worker pool [\#5](https://github.com/trento-project/runner/pull/5) (@arbulu89)
- Initial Web engine boilerplate [\#2](https://github.com/trento-project/runner/pull/2) (@arbulu89)

### Fixed

- Pass json list in quotes [\#49](https://github.com/trento-project/runner/pull/49) (@arbulu89)
- Set failed check executions as critical [\#43](https://github.com/trento-project/runner/pull/43) (@arbulu89)

### Closed Issues

- Remove unsued interval flag [\#21](https://github.com/trento-project/runner/issues/21)

### Other Changes

- remove check 1.3.7 [\#54](https://github.com/trento-project/runner/pull/54) (@gereonvey)
- Bump docker/metadata-action from 3.7.0 to 3.8.0 [\#52](https://github.com/trento-project/runner/pull/52) (@dependabot[bot])
- fix-1.1.1b: parser not to skip first line of the file [\#51](https://github.com/trento-project/runner/pull/51) (@fmherschel)
- fix-1.1.9: Corrected default for GCP \(copy and paste error\) [\#50](https://github.com/trento-project/runner/pull/50) (@fmherschel)
- Rename dev environment by default [\#48](https://github.com/trento-project/runner/pull/48) (@arbulu89)
- fix-1.3.3: platforms which does run the test \(no skip\) needs the exte… [\#47](https://github.com/trento-project/runner/pull/47) (@fmherschel)
- Fix 1.1.1 [\#45](https://github.com/trento-project/runner/pull/45) (@fmherschel)
- Remove interval flag [\#44](https://github.com/trento-project/runner/pull/44) (@arbulu89)
- Bump github.com/vektra/mockery/v2 from 2.12.0 to 2.12.1 [\#42](https://github.com/trento-project/runner/pull/42) (@dependabot[bot])
- fix 1.2.2: allow values to be greater-equal [\#41](https://github.com/trento-project/runner/pull/41) (@fmherschel)
- Adding skip variable for checks  [\#40](https://github.com/trento-project/runner/pull/40) (@pirat013)
- Bump github.com/vektra/mockery/v2 from 2.10.2 to 2.12.0 [\#39](https://github.com/trento-project/runner/pull/39) (@dependabot[bot])
- Fix 1.1.9b [\#38](https://github.com/trento-project/runner/pull/38) (@fmherschel)
- Fix 1.2.2 [\#37](https://github.com/trento-project/runner/pull/37) (@fmherschel)
- fix 1.3.4  [\#36](https://github.com/trento-project/runner/pull/36) (@fmherschel)
- Merge old project from  trento-project / trento  [\#34](https://github.com/trento-project/runner/pull/34) (@pirat013)
- change static value of msgwait to 2 times of watchdog [\#32](https://github.com/trento-project/runner/pull/32) (@schlosstom)
- Bump github.com/spf13/viper from 1.10.1 to 1.11.0 [\#31](https://github.com/trento-project/runner/pull/31) (@dependabot[bot])
- Check-rework for corosync and SBD [\#29](https://github.com/trento-project/runner/pull/29) (@pirat013)
- Bump actions/download-artifact from 2 to 3 [\#27](https://github.com/trento-project/runner/pull/27) (@dependabot[bot])
- Bump actions/upload-artifact from 2 to 3 [\#26](https://github.com/trento-project/runner/pull/26) (@dependabot[bot])
- Bump actions/setup-go from 2 to 3 [\#25](https://github.com/trento-project/runner/pull/25) (@dependabot[bot])
- Remove dev tag from ci/cd [\#23](https://github.com/trento-project/runner/pull/23) (@fabriziosestito)
- Change container tag to rolling [\#20](https://github.com/trento-project/runner/pull/20) (@fabriziosestito)
- Bump docker/metadata-action from 3.6.2 to 3.7.0 [\#19](https://github.com/trento-project/runner/pull/19) (@dependabot[bot])
- Bump github.com/vektra/mockery/v2 from 2.10.0 to 2.10.2 [\#17](https://github.com/trento-project/runner/pull/17) (@dependabot[bot])
- Bump github.com/swaggo/swag from 1.8.0 to 1.8.1 [\#15](https://github.com/trento-project/runner/pull/15) (@dependabot[bot])
- Remove unused api folder code [\#12](https://github.com/trento-project/runner/pull/12) (@arbulu89)
- Bump github.com/google/uuid from 1.1.2 to 1.3.0 [\#10](https://github.com/trento-project/runner/pull/10) (@dependabot[bot])
- Bump actions/cache from 2 to 3 [\#7](https://github.com/trento-project/runner/pull/7) (@dependabot[bot])
- Remove not used Start function and helpers from the Runner [\#6](https://github.com/trento-project/runner/pull/6) (@arbulu89)
- Remove duplicated ansible folder [\#4](https://github.com/trento-project/runner/pull/4) (@arbulu89)
- Build checks catalog [\#3](https://github.com/trento-project/runner/pull/3) (@arbulu89)
- Bump github.com/stretchr/testify from 1.7.0 to 1.7.1 [\#1](https://github.com/trento-project/runner/pull/1) (@dependabot[bot])

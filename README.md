# Helm tests generator

`helm-test-generator` looks into a folder (`test-scripts` by default) and generates a test job for each file it finds.

The job manifests and a secret containing the tests-scripts folder is written to stdout.

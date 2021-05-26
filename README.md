# Script to helm tests

`scripts-to-helm-test` looks into the folder `test-scripts` and generates a test job for each file it finds into `templates`.

For each test a secret is created to pass the script into the test job.

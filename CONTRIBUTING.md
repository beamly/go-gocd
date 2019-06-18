## How to contribute to `go-gocd`

### Reporting Bugs

 - **Ensure the bug was not already reported**by searching on GitHub under [Issues](https://github.com/beamly/go-gocd/issues).
 - If you're unable to find an open issue addressing the problem, [open a new one](https://github.com/beamly/go-gocd/issues/new). Be sure to include a **title and clear description**, as much relevant information as possible, and a **code sample** or an **executable test case** demonstrating the expected behavior that is not occurring.

### Fixing Bugs

- Open a new GitHub pull request with the patch.
- Ensure the PR description clearly describes the problem and solution. Include the relevant issue number if applicable.

### Adding Features

 - Right now, the best feature you can implement is one of the API resources. If you would like to add a new endpoint for the API, have a look at the [Encryption endpoint](https://github.com/beamly/go-gocd/blob/master/gocd/encryption.go) as a basic example of the structure of an API action.
 - Make sure to run `make format` and `make test` before submitting a PR.

## Running Tests

    git clone https://github.com/beamly/go-gocd.git
    cd go-gocd
    make test

## Opening PR's

 - Make sure all tests pass in travis
 - Add test cases to any code you change
 - Update the library readme by running `make doc`

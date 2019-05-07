# Form3 payments api [![Build Status](https://travis-ci.org/florinutz/form3-api.svg?branch=master)](https://travis-ci.org/florinutz/form3-api)

* You should use best practice, for example TDD/BDD, with a focus on full-stack testing
* Prioritize correctness, robustness, and extensibility over extra features and optimizations.
* Write your code with the quality bar you would use for production code.
* Try to simplify your code by using open source frameworks and libraries where possible

## Mongo setup

Start:
```bash
docker run -d --name mongo \
    -p 27017-27019:27017-27019 mongo \
    --bind_ip_all
```

## Usage
```bash
make binary
MONGO_DSN=localhost ./bin
```
Use make and check out the tests!
Payment routes are mounted under the /payments namespace

### Generate router docs
(quite ugly)
```bash
./bin -doc > doc_filename.md
```
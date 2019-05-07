# Form3 payments api

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

And don't forget the env DSN:

```bash
MONGO_DSN=localhost
```

## Usage

Use make and check out the tests!
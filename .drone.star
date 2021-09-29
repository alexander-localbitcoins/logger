def main(ctx):
    return [
        {
            "kind": "pipeline",
            "type": "docker",
            "name": "dronesecret",
            "trigger": {"event": ["push"]},
            "steps": [
                test_step,
            ],
            "node": {"docker": "slow"},
        },
    ]

test_step = {
    "name": "tests",
    "image": "golang:latest",
    "commands": [
        "go test -v",
        "cd ./mock/",
        "go test -v",
    ],
}

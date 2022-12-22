load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("//:gqlgen.bzl", "gqlgen")

go_library(
    name = "graph",
    deps = [":gqlgen"],
)

gqlgen(
    name = "gqlgen",
    schemas = glob(["schema/**/*.gql"]),
)

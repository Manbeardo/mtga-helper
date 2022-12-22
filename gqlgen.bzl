load("@bazel_skylib//lib:paths.bzl", "paths")
load(
    "@io_bazel_rules_go//go:def.bzl",
    _go_context = "go_context",
)

_CONFIG = """
schema:
  - {schema}
exec:
  filename: {exec_path}
  package: {go_package}
model:
  filename: {model_path}
  package: {go_package}
resolver:
  layout: follow-schema
  dir: .
  package: {go_package}
  filename_template: {resolver_path}
"""

def _gqlgen_impl(ctx):
    go = _go_context(ctx)

    config_file = ctx.actions.declare_file(ctx.label.name + ".yml")
    config_content = _CONFIG.format(
        schema = ctx.attr.schema.label.name,
        go_package = ctx.attr.go_package,
        exec_path = ctx.outputs.exec_file.basename,
        model_path = ctx.outputs.model_file.basename,
        resolver_path = ctx.outputs.resolver_file.basename,
    )
    print(config_content)
    ctx.actions.write(
        content = config_content,
        output = config_file,
    )

    out_files = [ctx.outputs.exec_file, ctx.outputs.model_file, ctx.outputs.resolver_file]
    inputs = [config_file]

    print(go.env)

    env = {
        "GOCACHE": "/tmp",
        "GOPACKAGESPRINTGOLISTERRORS": "1",
        "GOPACKAGESDEBUG": "1",
        "GOPACKAGESDRIVER": ctx.executable._gopackagesdriver.path,
        "GODEBUG": "execerrdot=0",
    }

    env.update(go.env)

    path_elements = go.env["PATH"].split(ctx.configuration.host_path_separator)
    path_elements.append("{}/bin".format(go.env["GOROOT"]))
    path_elements.append(ctx.executable._gopackagesdriver.dirname)

    env["PATH"] = ctx.configuration.host_path_separator.join(path_elements)

    print(env)
    print(ctx.bin_dir.path)
    print(ctx.executable._gqlgen.path)

    ctx.actions.run(
        tools = [go.go, ctx.executable._gopackagesdriver],
        inputs = inputs,
        outputs = out_files,
        arguments = ["--config", config_file.path, "--verbose"],
        executable = ctx.executable._gqlgen,
        env = env,
    )

_gqlgen = rule(
    implementation = _gqlgen_impl,
    attrs = {
        "go_package": attr.string(
            doc = "Go package name for the generated files",
            mandatory = True,
        ),
        "exec_file": attr.output(
            doc = "File for the generated server implementation",
            mandatory = True,
        ),
        "model_file": attr.output(
            doc = "File for the generated model spec",
            mandatory = True,
        ),
        "resolver_file": attr.output(
            doc = "File for the partially-generated resolver implementation",
            mandatory = True,
        ),
        "schema": attr.label(
            allow_files = [".graphqls", ".graphql", ".gql"],
            doc = "The schema file from which to generate Go code",
            mandatory = True,
        ),
        "_gqlgen": attr.label(
            default = "@com_github_99designs_gqlgen//:gqlgen",
            cfg = "exec",
            executable = True,
        ),
        "_gopackagesdriver": attr.label(
            default = "@io_bazel_rules_go//go/tools/gopackagesdriver:gopackagesdriver",
            cfg = "exec",
            executable = True,
        ),
        "_go_context_data": attr.label(
            default = "@io_bazel_rules_go//:go_context_data",
        ),
    },
    toolchains = ["@io_bazel_rules_go//go:toolchain"],
)

def gqlgen(name, schemas):
    concatenated_schemas = ":complete_schema.gql"
    native.genrule(
        name = name + "_concat_schemas",
        srcs = schemas,
        outs = [concatenated_schemas],
        cmd = "cat $(OUTS) > $@",
    )

    _gqlgen(
        name = name,
        go_package = paths.basename(native.package_name()),
        schema = concatenated_schemas,
        model_file = ":models.go",
        exec_file = ":exec.go",
        resolver_file = ":resolvers.go",
    )

**This is a third-party repository managed by Buf.**

This module contains the core types that almost all developers use from the
[googleapis](https://github.com/googleapis/googleapis) repository, specifically:

- The core files from the `google.api` package, used by many for the HTTP annotations including
  grpc-gateway.
- The packages `google.api.expr.v1alpha1`, `google.api.expr.v1beta1`, used within Envoy.
- The `google.bytestream` package.
- The `google.longrunning` package.
- The `google.geo.type` package.
- The `google.rpc` package, used by grpc.
- The `google.type` package, containing many common types.

The [source repository](https://github.com/googleapis/googleapis) contains over 3800 files, mostly
relating to Google's core APIs. However, the ~30 files above are the only files used by 99.999% of
developers, and these files are the most common dependency in the Protobuf ecosystem. This hosted
module only includes these specific files, as including all the files causes hundreds of megabytes
of unused generated code for the vast majority of developers. To use Google's core APIs, create your
own module that has a `dep` on `buf.build/googleapis/googleapis` with the specific packages you want
to use.

Updates to the [source repository](https://github.com/googleapis/googleapis) are automatically
synced on a periodic basis, and each BSR commit is tagged with corresponding Git commits.

To depend on a specific Git commit, you can use it as your reference in your dependencies:

```
deps:
  - buf.build/googleapis/googleapis:<GIT_COMMIT_TAG>
```

For more information, see the [documentation](https://docs.buf.build/bsr/overview).

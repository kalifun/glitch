# Examples

This directory contains sample YAML definitions and their generated Go output.

## Layout
- `examples/errors/`: input YAML files, split by domain.
- `examples/generated/`: generated Go code (do not edit manually).

## Generate
From the repo root:

```bash
# Generate into examples/generated using package name "generated"
go run . gen -y examples/errors/*.yaml -p generated --out examples/generated
```

## Notes
- Output filenames follow the YAML filenames, e.g. `errors/auth.yaml` -> `generated/auth.go`.
- The generated package name is controlled by `-p` (here `generated`).

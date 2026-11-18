#!/usr/bin/env python3
import json
import pathlib
import subprocess
import unittest

ROOT = pathlib.Path(__file__).resolve().parents[2]
SCRIPT = ROOT / "scripts/release/generate_artifact_manifest.py"


class ArtifactManifestTests(unittest.TestCase):
    def test_generates_neutral_mingligo_binding(self):
        result = subprocess.run(
            [
                "python3", str(SCRIPT), "--repo", str(ROOT),
                "--image", f"r/mingligo-mcp@sha256:{'a'*64}",
            ],
            check=True,
            capture_output=True,
            text=True,
        )
        document = json.loads(result.stdout)
        self.assertEqual(set(document), {"kind", "source_commit", "image"})
        self.assertEqual(document["kind"], "mingligo-artifact")
        self.assertNotIn("hk.liki", result.stdout)

    def test_rejects_mutable_image(self):
        result = subprocess.run(
            ["python3", str(SCRIPT), "--repo", str(ROOT), "--image", "r/mingligo-mcp:latest"],
            capture_output=True,
            text=True,
        )
        self.assertNotEqual(result.returncode, 0)


if __name__ == "__main__":
    unittest.main()

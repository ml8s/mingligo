#!/usr/bin/env python3
import json
import pathlib
import subprocess
import unittest

ROOT = pathlib.Path(__file__).resolve().parents[2]
SCRIPT = ROOT / "scripts/release/generate_capability_manifest.py"


class CapabilityManifestTests(unittest.TestCase):
    def test_generates_neutral_engine_capability(self):
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
        self.assertEqual(document["kind"], "mingligo-capability")
        self.assertEqual(document["project"], "mingligo")
        self.assertNotIn("counsel", document["images"])
        self.assertNotIn("hk.liki", result.stdout)
        self.assertEqual(len(document["tools"]["engine"]), 17)

    def test_rejects_mutable_image(self):
        result = subprocess.run(
            ["python3", str(SCRIPT), "--repo", str(ROOT), "--image", "r/mingligo-mcp:latest"],
            capture_output=True,
            text=True,
        )
        self.assertNotEqual(result.returncode, 0)


if __name__ == "__main__":
    unittest.main()

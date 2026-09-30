"""Custom unittest metadata. The decorator preserves the original function."""
import unittest


def cliewen(*, ac, type, direction):
    def annotate(method):
        method.cliewen = {"ac": ac, "type": type, "direction": direction}
        return method
    return annotate


class ExampleTest(unittest.TestCase):
    @cliewen(ac="AC-001", type="Unit", direction="positive")
    def test_valid_input(self):
        self.assertEqual(int("42"), 42)

    @cliewen(ac="AC-001", type="Unit", direction="negative")
    def test_invalid_input(self):
        with self.assertRaises(ValueError):
            int("not an integer")

import sys
from pathlib import Path

PROTO_DIR = str(Path(__file__).parent)
if PROTO_DIR not in sys.path:
    sys.path.append(PROTO_DIR)


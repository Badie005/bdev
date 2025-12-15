# Python CLI Template

Template B.DEV pour créer des CLI Python avec Typer.

## Stack

- Python 3.10+
- Typer (CLI)
- Rich (affichage)
- pytest (tests)

## Usage

```bash
bdev new python-cli mon-cli
cd mon-cli
python -m venv .venv
.venv\Scripts\activate
pip install -r requirements.txt
python main.py --help
```

## Structure

```
├── main.py           # Point d'entrée
├── commands/         # Commandes CLI
├── utils/            # Utilitaires
├── tests/            # Tests
└── requirements.txt
```

## Commandes

```bash
python main.py          # Exécuter
python main.py --help   # Aide
pytest                  # Tests
```

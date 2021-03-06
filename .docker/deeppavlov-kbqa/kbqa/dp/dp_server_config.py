from pathlib import Path

from deeppavlov.utils import settings
from deeppavlov.core.common.file import read_json, save_json

settings_path = Path(settings.__path__[0]) / 'server_config.json'

settings = read_json(settings_path)

settings['model_defaults'] = settings.get('model_defaults', {})


settings['model_defaults']['kbqa'] = {
    "model_endpoint": "/answers",
    "model_args_names": ["sentences"]
}

print(settings_path)
save_json(settings, settings_path)


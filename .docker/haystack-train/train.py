#!/usr/bin/env python3

HOST = 'elastic'
PORT = 9200
INDEX_NAME = 'emlyon'

from haystack import Finder
from haystack.preprocessor.cleaning import clean_wiki_text
from haystack.preprocessor.utils import convert_files_to_dicts, fetch_archive_from_http
from haystack.reader.farm import FARMReader
from haystack.reader.transformers import TransformersReader
from haystack.utils import print_answers
from haystack.document_store.elasticsearch import ElasticsearchDocumentStore

document_store = ElasticsearchDocumentStore(host=HOST, port=PORT, username="", password="", index=INDEX_NAME)

# clear existing index (optional)
if document_store.client.indices.exists(index=document_store.index):
    print('clear existing inddex')
    document_store.client.indices.delete(index=document_store.index)

# Get all dirs in wikipedia folder
from os import listdir
from os.path import isfile, join
import json
from tqdm import tqdm

emlyon_path = "/opt/em-lyon/data"
onlydirs = [f for f in listdir(emlyon_path) if not isfile(join(emlyon_path, f))]

dicts = []
bulk_size = 5000

pbar = tqdm(onlydirs)
for directory in pbar:
    subdirs = [f for f in listdir(join(emlyon_path,directory)) if not isfile(join(emlyon_path,directory))]
    pbar.set_description(f"Processing emlyon folder {directory}")

    for file in subdirs:
        f = open(join(emlyon_path, directory, file), "r")

        # Each text file contains json structures separated by EOL
        articles = f.read().split("\n")

        for article in articles:
            if len(article)==0: continue

            # Article in json format
            json_formatted_article = json.loads(article)

            # Rename keys
            document = {"id": json_formatted_article["hash"],
                        "name": json_formatted_article["title"],
                        "url": json_formatted_article["url"],
                        "text": json_formatted_article["body"],
                        "canonical-link": json_formatted_article["canonical-link"],
                        "image": json_formatted_article["image"],
                        "publish-date": json_formatted_article["publish-date"],
                        "meta-description": json_formatted_article["description"],
                        "meta-keywords": json_formatted_article["keywords"],
                        "meta-lang": json_formatted_article["lang"],
                        "language": json_formatted_article["language"]}


            # Add document to bulk
            dicts.append(document)

            if len(dicts)>=bulk_size:
                # Index bulk
                try:
                    document_store.write_documents(dicts)
                    dicts.clear()
                except:
                    print("Bulk not indexed")


if len(dicts) > 0:
    print('final round')
    document_store.write_documents(dicts)

print('finished')

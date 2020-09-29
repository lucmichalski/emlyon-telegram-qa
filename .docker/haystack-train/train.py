#!/usr/bin/env python3

from haystack import Finder
from haystack.preprocessor.cleaning import clean_wiki_text
from haystack.preprocessor.utils import convert_files_to_dicts, fetch_archive_from_http
from haystack.reader.farm import FARMReader
from haystack.reader.transformers import TransformersReader
from haystack.utils import print_answers
from haystack.document_store.elasticsearch import ElasticsearchDocumentStore

# Get all dirs in wikipedia folder
from os import listdir
from os.path import isfile, join
import json
from tqdm import tqdm

import click

def train_model(data_dir, es_host, es_port, index_name, index_reset, es_username, es_password):

    document_store = ElasticsearchDocumentStore(host=es_host, port=es_port, username=es_username, password=es_password, index=index_name)

    # clear existing index (optional)
    if index_reset:
        if document_store.client.indices.exists(index=document_store.index):
            print('clear existing inddex')
            document_store.client.indices.delete(index=document_store.index)

    # onlydirs = [f for f in listdir(data_dir) if not isfile(join(data_dir, f))]
    onlydirs = [data_dir]

    dicts = []
    bulk_size = 500

    pbar = tqdm(onlydirs)
    for directory in pbar:
        subdirs = [f for f in listdir(join(data_dir,directory)) if not isfile(join(data_dir,directory))]
        pbar.set_description(f"Processing folder {directory}")

        for file in subdirs:
            f = open(join(data_dir, directory, file), "r")

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
                            "language": json_formatted_article["detected-lang"],
			    "fingerprint": json_formatted_article["fingerprint"]}

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

@click.command()
@click.option("--data-dir", default="/opt/em-lyon/data/en", help="Data directory.")
@click.option("--host", default="elastic", help="Elastic host.")
@click.option("--port", default="9200", help="Elastic port.")
@click.option("--username", default="",  help="Elastic username.")
@click.option("--password", default="",  help="Elastic password.")
@click.option("--index-name", default="emlyon-en",  help="Elastic index name (nb. must be lowercase).")
@click.option("--index-reset", default=False, is_flag=True, help="Elastic index truncate.")
def service(data_dir, host, port, index_name, index_reset, username, password):
    train_model(data_dir, host, port, index_name, index_reset, username, password)

if __name__ == '__main__':
    service()

#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Fri Sept 11 13:23:10 2020
@author: luc michalski
"""
import os
import json

import logging
from logging.handlers import RotatingFileHandler

import argparse

import tqdm
import numpy as np  # linear algebra
import pandas as pd # data processing, CSV file I/O (e.g. pd.read_csv)

from flask import Flask, jsonify, request

from haystack import Finder
from haystack.preprocessor.cleaning import clean_wiki_text
from haystack.preprocessor.utils import convert_files_to_dicts, fetch_archive_from_http
from haystack.reader.farm import FARMReader
from haystack.reader.transformers import TransformersReader
from haystack.utils import print_answers

from haystack.document_store.elasticsearch import ElasticsearchDocumentStore
from haystack.retriever.sparse import ElasticsearchRetriever

# Script arguments can include path of the config
arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('--model', type=str, default="deepset/roberta-base-squad2")
arg_parser.add_argument('--host', type=str, default="0.0.0.0")
arg_parser.add_argument('--port', type=str, default="8000")
arg_parser.add_argument('--log', type=str, default="../logs/ovh-qa.log")
args = arg_parser.parse_args()

def filter_answers(results: dict, details: str = "all"):
    answers = results["answers"]
    if details != "all":
        if details == "minimal":
            keys_to_keep = set(["answer", "context"])
        elif details == "medium":
            keys_to_keep = set(["answer", "context", "score"])
        else:
            keys_to_keep = answers.keys()

        # filter the results
        filtered_answers = []
        for ans in answers:
            filtered_answers.append({k: ans[k] for k in keys_to_keep})
        return filtered_answers
    else:
        return results


app = Flask(__name__)
reader = FARMReader(model_name_or_path=args.model, use_gpu=True, context_window_size=500)

@app.route('/query')
def query():

    if request.args.get('question'):
        question = request.args.get('question')
    else:
        result = {"status": 400, "msg": "Question cannot be empty"}
        return jsonify(result)

    if request.args.get('index'):
        index = request.args.get('index')
    else:
        index = "ovh-en-us"

    if request.args.get('top_k_reader'):
        top_k_reader = request.args.get('top_k_reader')
    else:
        top_k_reader = 20

    if request.args.get('top_k_retriever'):
        top_k_retriever = request.args.get('top_k_retriever')
    else:
        top_k_retriever = 2

    document_store = ElasticsearchDocumentStore(host="elastic", username="", password="", index=index)
    retriever = ElasticsearchRetriever(document_store=document_store)
    finder = Finder(reader, retriever)
    prediction = finder.get_answers(question=question, top_k_retriever=int(top_k_retriever), top_k_reader=int(top_k_reader))
    result = filter_answers(prediction, details="all")
    app.logger.info('question: %s', question)
    app.logger.info('result: %s', result)
    return jsonify(result)

if __name__ == '__main__':
    handler = RotatingFileHandler(args.log, maxBytes=10000, backupCount=1)
    handler.setLevel(logging.INFO)
    app.logger.addHandler(handler)
    app.run(host=args.host, port=args.port)

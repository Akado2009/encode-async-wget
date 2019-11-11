import pandas as pd

import uuid

import argparse
import datetime

parser = argparse.ArgumentParser(description='Uploader into clickhouse')
parser.add_argument('-f', action='store',
                    dest='first_accession',
                    help='First accession number')
parser.add_argument('-s', action='store',
                    dest='second_accession',
                    help='Second accession number')
parser.add_argument('-r', action='store',
                    dest='result_folder',
                    help='Result folder')
args = parser.parse_args()

#Generate ID

first_id = str(uuid.uuid4())
second_id = str(uuid.uuid4())

#Background

bkg_df = pd.read_csv('{}-{}.bkg', format(
    args.first_accession, args.second_accession
), header=None)
bkg_df.insert(0, "accession2", args.second_accession)
bkg_df.insert(0, "accession1", args.first_accession)
bkg_df[-1] = datetime.datetime.now().strftime("%Y-%m-%d")
bkg_df.to_csv('{}-{}_editted.bkg', format(
    args.first_accession, args.second_accession
), header=False, quotechar="'", index=False)

#Distribution

dist_df = pd.read_csv('{}-{}.dist', format(
    args.first_accession, args.second_accession
), header=None, skiprows=2, sep='\t')
dist_df.insert(0, "accession2", args.second_accession)
dist_df.insert(0, "accession1", args.first_accession)
dist_df[-1] = datetime.datetime.now().strftime("%Y-%m-%d")
dist_df.to_csv('{}-{}_editted.dist', format(
    args.first_accession, args.second_accession
), header=False, quotechar="'", index=False)
#Foreground

fg_df = pd.read_csv('{}-{}.fg', format(
    args.first_accession, args.second_accession
), header=None, sep='\t')
fg_df.insert(0, "accession2", args.second_accession)
fg_df.insert(0, "accession1", args.first_accession)
fg_df[-1] = datetime.datetime.now().strftime("%Y-%m-%d")
fg_df.to_csv('{}-{}_editted.fg', format(
    args.first_accession, args.second_accession
), header=False, quotechar="'", index=False)


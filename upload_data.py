import pandas as pd

import uuid

import argparse
import datetime
# tissue           │ String │              │                    │         │                  │
# │ cell_line        │ String │              │                    │         │                  │
# │ dataset          │ String │              │                    │         │                  │
# │ feature 
parser = argparse.ArgumentParser(description='Uploader into clickhouse')
parser.add_argument('-f', action='store',
                    dest='first_accession',
                    help='First accession number')
parser.add_argument('--fp', action='store',
                    dest='first_parameters',
                    help='First accession parameters')
parser.add_argument('-s', action='store',
                    dest='second_accession',
                    help='Second accession number')
parser.add_argument('--sp', action='store',
                    dest='second_parameter',
                    help='Second accession parameters')
parser.add_argument('-r', action='store',
                    dest='result_folder',
                    help='Result folder')
args = parser.parse_args()

#Generate ID

first_id = str(uuid.uuid4())
second_id = str(uuid.uuid4())

#General_info
#Statistics
#Parameters
#Background

bkg_df = pd.read_csv('{}-{}.bkg', format(
    first_id, second_id
), header=None)
bkg_df.insert(0, "accession2", second_id)
bkg_df.insert(0, "accession1", first_id)
bkg_df[-1] = datetime.datetime.now().strftime("%Y-%m-%d")
bkg_df.to_csv('{}-{}_editted.bkg', format(
    first_id, second_id
), header=False, quotechar="'", index=False)

#Distribution

dist_df = pd.read_csv('{}-{}.dist', format(
    first_id, second_id
), header=None, skiprows=2, sep='\t')
dist_df.insert(0, "accession2", second_id)
dist_df.insert(0, "accession1", first_id)
dist_df[-1] = datetime.datetime.now().strftime("%Y-%m-%d")
dist_df.to_csv('{}-{}_editted.dist', format(
    first_id, second_id
), header=False, quotechar="'", index=False)


#Foreground

fg_df = pd.read_csv('{}-{}.fg', format(
    first_id, second_id
), header=None, sep='\t')
fg_df.insert(0, "accession2", second_id)
fg_df.insert(0, "accession1", first_id)
fg_df[-1] = datetime.datetime.now().strftime("%Y-%m-%d")
fg_df.to_csv('{}-{}_editted.fg', format(
    first_id, second_id
), header=False, quotechar="'", index=False)


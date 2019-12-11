import pandas as pd
import os
import uuid
import datetime

# Iterate over table
# Fill general_info

# CONSTS
ROOT_DIRECTORY = "/mnt/scratch/shared/SG_KIRILL/results"
GENERAL_TABLE = "encode.files.txt"


def get_folders(root):
    return os.listdir(root)

def parse_general(table_name, out_file, mapping):
    df = pd.read_csv(table_name, sep='\t')
    for _, row in df.iterrows():
        unique_id = uuid.uuid4()
        mapping[row["Accession"]] = unique_id
        cell_line = row["Tissue"]
        if cell_line == ".":
            cell_line = row["CellLine"]
        if cell_line == ".":
            cell_line = row["PrimaryCell"]
        dataset = row["Dataset"].split("/")[2]
        out_file.write(
            "{},{},{},{},{},{},{}\n".format(
                unique_id, row["Accession"], "empty",
                cell_line, dataset, row["Feature"],
                datetime.datetime.now().strftime("%Y-%m-%d")
            )
        )

if __name__ == "__main__":
    general_info_file = open("general_info.tab", "w")
    id_to_uuid = {}
    parse_general(GENERAL_TABLE, general_info_file, id_to_uuid)
    general_info_file.close()
    
    results = get_folders(ROOT_DIRECTORY)
    # for folder in results:
    print(id_to_uuid)
    print(len(results))
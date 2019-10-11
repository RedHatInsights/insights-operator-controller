#!/bin/sh

DATABASE=controller.db

cat schema.sql | sqlite3 ../${DATABASE}
cat test_data.sql | sqlite3 ../${DATABASE}


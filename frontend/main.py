import ast
import os
from datetime import datetime

import altair as alt
import matplotlib.pyplot as plt
import pandas as pd
import psycopg2
import streamlit as st

# Set up database connection
DB_HOST = os.environ.get('POSTGRES_HOST')
DB_PORT = os.environ.get('POSTGRES_PORT')
DB_NAME = os.environ.get('POSTGRES_DB')
DB_USER = os.environ.get('POSTGRES_USER')
DB_PASSWORD = os.environ.get('POSTGRES_PASSWORD')

with open("/bindings/tanzutrends-db/dbname") as f:
    DB_NAME = f.read()
with open("/bindings/tanzutrends-db/instancename") as f:
    DB_HOST = f.read()
with open("/bindings/tanzutrends-db/username") as f:
    DB_USER = f.read()
with open("/bindings/tanzutrends-db/password") as f:
    DB_PASSWORD = f.read()
# Set port manual
DB_PORT = 5432

if not all([DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD]):
    st.error('One or more required environment variables are missing.')
    st.stop()

DB_URL = f"postgresql://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}"
conn = psycopg2.connect(DB_URL)
cur = conn.cursor()
cur.execute("SELECT * FROM tweets")
rows = cur.fetchall()

# Convert list of tuples to Pandas DataFrame
df = pd.DataFrame(rows, columns=[desc[0] for desc in cur.description])

# Close database cursor and connection
cur.close()
conn.close()

# Replace all names with lowercase
df['hashtags'] = df['hashtags'].str.lower()

df['time'] = pd.to_datetime(df['time'])


# Split hashtags into separate objects
df['hashtags'] = df['hashtags'].str.split(',')

# Create a new dataframe with each hashtag and its count
hashtags_counts_df = pd.DataFrame({
    'hashtags': [hashtag for hashtags in df['hashtags'] for hashtag in hashtags]
})
hashtags_counts_df = hashtags_counts_df['hashtags'].value_counts().reset_index()
hashtags_counts_df.columns = ['hashtags', 'count']

# Create a multi-select dropdown for selecting hashtags to plot
selected_hashtags = st.sidebar.multiselect(
    'Select hashtags to plot',
    options=hashtags_counts_df['hashtags'].unique(),
    default=hashtags_counts_df['hashtags'].unique()[:3]
)

# Filter the data based on the selected hashtags
filtered_df = df[df['hashtags'].apply(lambda hashtags: any(hashtag in selected_hashtags for hashtag in hashtags))]

# Convert the time column to a datetime object
filtered_df["time"] = pd.to_datetime(filtered_df["time"])

# Group the data by day and count the number of occurrences
grouped_df = filtered_df.groupby(filtered_df["time"].dt.date)["hashtags"].count().reset_index()

# Create a chart for each selected hashtag
charts = []
for hashtag in selected_hashtags:
    hashtag_filtered_df = filtered_df[filtered_df['hashtags'].apply(lambda hashtags: hashtag in hashtags)]
    hashtag_grouped_df = hashtag_filtered_df.groupby(filtered_df["time"].dt.date)["hashtags"].count().reset_index()

    chart = alt.Chart(hashtag_grouped_df).mark_line().encode(
        x=alt.X("time:T", title="Date"),
        y=alt.Y("hashtags:Q", title="Number of occurrences"),
        tooltip=["time", "hashtags"]
    ).properties(
        width=800,
        height=400,
        title=hashtag
    ).interactive()

    charts.append(chart)

st.header("Tanzu Trends")


# Combine the charts into a single row
st.altair_chart(alt.hconcat(*charts), use_container_width=True)
# hybrid-search-demo


To get started you need to create an Elasticsearch instace. It is built around serverless.

Here is the documentation for setting up and reindexing the indexes
https://www.elastic.co/docs/solutions/search/hybrid-semantic-text


You will need to import data.json into your data ingest index and then reindex into your semantic_text index.


Data is from here as of 2025-01-22:
https://open.canada.ca/data/en/dataset/d38de914-c94c-429b-8ab1-8776c31643e3

This is the semantic index mapping for this data
```
PUT your_index_name
{
  "mappings": {
    "properties": {
      "semantic_text": {
        "type": "semantic_text"
      },
      "Product": {
        "type": "text",
        "copy_to": "semantic_text"
      }
    }
  }
}
```

You will need to setup 3 environmental variables:
- ES_API_KEY = Your ES API Key
- ES_SERVER_URL = full URL to your severless Elasticsearch Instance (https://????:443)
- ES_SEARCH_INDEX = the name of the index you created for the semantic_text using the mapping above

Then run the app with no cli flags. This will launch a Terminal User Interface (TUI).

```
$ go run .
```

### Usage
- **Type** your search query into the input box.
- Press **Enter** to execute the search.
- Use the **Up/Down Arrow Keys** to scroll through results.
- Press **Ctrl+C** or **Esc** to quit the application.

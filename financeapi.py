import yfinance as yf
from flask import Flask, jsonify
import pandas as pd

app = Flask(__name__)

@app.route('/data', methods=['GET'])
def get_data():
    print("recieved request")
    data = yf.download("SPY", period="5d", interval='1h')
    data.reset_index(inplace=True)

    print(data.head())
    data['Datetime'] = data['Datetime'].astype(str)
    return jsonify(data.to_dict())

if __name__ == '__main__':
    app.run(port=5000)

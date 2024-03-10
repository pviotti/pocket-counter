# Module to
# 1. serve static HTML page
# 2. serve API to retrieve unread counts
# 3. periodically query Pocket API and retrieve unread count

from flask import Flask, request, send_from_directory

# app = Flask(__name__, static_url_path='', static_folder='static')
# db = Database()

# if app.debug:
#     # pip install flask-cors
#     from flask_cors import CORS
#     CORS(app)

# @app.route("/grep")
# def grep():
#     term = request.args.get('q')
#     results = db.grep(term)
#     return [res.to_dict() for res in results]

# @app.route("/", defaults={'path':''})
# def root(path):
#     return send_from_directory(app.static_folder,'index.html')



import requests
import json
import time
import webbrowser

# Replace with your actual Pocket API consumer key
CONSUMER_KEY = '108694-ac09bf3184e4d49360df137'
redirect_uri = "https://www.ilpost.it/"

def get_unread_articles_count(access_token):
    url = 'https://getpocket.com/v3/get'
    headers = {'Content-Type': 'application/json; charset=UTF-8',
               'X-Accept': 'application/json'}
    data = {
        'consumer_key': CONSUMER_KEY,
        'access_token': access_token,
        'state': 'unread',
    }

    response = requests.post(url, headers=headers, data=json.dumps(data))
    if response.status_code == 200:
        response_data = json.loads(response.content)
        unread_count = len(response_data['list'])
        return unread_count
    else:
        print(f"Error: {response.status_code} - {response.content}")
        return None

def get_request_token():
    url = 'https://getpocket.com/v3/oauth/request'
    data = {'consumer_key': CONSUMER_KEY, 'redirect_uri': redirect_uri}
    response = requests.post(url, json=data)
    request_token = response.text.split('=')[1]
    return request_token

def authorize_request_token(request_token):
    auth_url = f"https://getpocket.com/auth/authorize?request_token={request_token}&redirect_uri={redirect_uri}"
    webbrowser.open(auth_url)

def get_access_token(authorized_request_token):
    url = 'https://getpocket.com/v3/oauth/authorize'
    data = {'consumer_key': CONSUMER_KEY, 'code': authorized_request_token}
    response = requests.post(url, json=data)
    print(response.text)
    access_token = response.text.split('&')[0].split('=')[1]
    return access_token


def main():
    #while True:
    unread_count = get_unread_articles_count()
    if unread_count is not None:
        # Store the unread_count in your local data store (e.g., a file or database)
        print(f"Unread articles: {unread_count}")
    #time.sleep(24 * 60 * 60)  # Wait for 24 hours before checking again

if __name__ == '__main__':
#    request_token = get_request_token()
#    print(f"request_token: {request_token}")
#    authorize_request_token(request_token)
#
#    time.sleep(15)
#    access_token = get_access_token(request_token)
#    print(f"access_token: {access_token}")
    access_token = 'e9902fdb-ac38-f259-dc9b-da9a56'
    unread_count = get_unread_articles_count(access_token)
    print(f"unread count: {unread_count}")


    #main()
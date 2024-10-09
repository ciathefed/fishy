import requests
from bs4 import BeautifulSoup


url = 'https://chromium.googlesource.com/chromiumos/docs/+/master/constants/errnos.md'

response = requests.get(url)

soup = BeautifulSoup(response.content, 'html.parser')

table = soup.find('table')

results = []

for row in table.find_all('tr'):
    cells = row.find_all('td')
    
    if len(cells) >= 4:
        third_item = cells[2].get_text(strip=True)  
        fourth_item = cells[3].get_text(strip=True)
        
        results.append((third_item, fourth_item))

for item in results:
    # print(f'{item[0]}: "{item[1].lower()}",')
    print(f'{item[0]},')

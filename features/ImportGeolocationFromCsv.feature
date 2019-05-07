Feature: ImportGeolocationFromCsv
  As a user
  I want to import geolocation from csv file
  so that I save geolocation details


  Scenario: Import geolocation from csv file
    Given there is a csv file in the path "/tmp/data.csv" with the following content
    """
    ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
    """
    When I run a command "geolocation-service-import-data" with args "-f /tmp/data.csv"
    Then there should be "3" geolocation(s) stored in the table "geolocation"
    And the following geolocation(s) should be stored in the table "geolocation"
      | ip_address     | country_code | country      | city         | latitude           | longitude          | mystery_value |
      | 200.106.141.15 | SI           | Nepal        | DuBuquemouth | -84.87503094689836 | 7.206435933364332  | 7823011346    |
      | 160.103.7.140  | CZ           | Nicaragua    | New Neva     | -68.31023296602508 | -37.62435199624531 | 7301823115    |
      | 70.95.73.73    | TL           | Saudi Arabia | Gradymouth   | -49.16675918861615 | -86.05920084416894 | 2559997162    |

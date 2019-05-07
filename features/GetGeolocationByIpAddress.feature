Feature: GetGeolocationByIpAddress
  As a user
  I want to know the geolocation details
  from a given ip

  Background:
    Given that the following geolocation(s) are stored in the table "geolocation"

      | id                                   | ip_address     | country_code | country      | city         | latitude           | longitude          | mystery_value |
      | 8b3f0880-800c-40cf-9cc6-2d53be233c3f | 200.106.141.15 | SI           | Nepal        | DuBuquemouth | -84.87503094689836 | 7.206435933364332  | 7823011346    |
      | 5c6a6ae0-d005-4123-bfb0-4758594ae3b8 | 160.103.7.140  | CZ           | Nicaragua    | New Neva     | -68.31023296602508 | -37.62435199624531 | 7301823115    |
      | 71aa6f04-ede9-46f4-a63d-373c3c206fc1 | 70.95.73.73    | TL           | Saudi Arabia | Gradymouth   | -49.16675918861615 | -86.05920084416894 | 2559997162    |


  Scenario: Geolocation details should be displayed
    When I request REST endpoint with method "GET" and path "/geolocation/160.103.7.140"

    Then I should have an OK response

    And I should have a response with following JSON body
    """
    {
      "id": "5c6a6ae0-d005-4123-bfb0-4758594ae3b8",
      "ip_address": "160.103.7.140",
      "city": "New Neva",
      "country": "Nicaragua",
      "country_code": "CZ",
      "id": "5c6a6ae0-d005-4123-bfb0-4758594ae3b8",
      "ip_address": "160.103.7.140",
      "latitude": "-68.31023296602508",
      "longitude": "-37.62435199624531",
      "mystery_value": 7301823115
    }
    """

  Scenario: Geolocation details should not be displayed, ip address does not exists
    When I request REST endpoint with method "GET" and path "/geolocation/160.103.7.145"

    Then I should have a not found response

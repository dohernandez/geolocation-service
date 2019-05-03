Feature: REST Context

  Scenario: Successful GET Request
    When I request REST endpoint with method "GET" and path "/get-something?foo=bar"

    Then I should have an OK response

    And I should have a response with following JSON body
    """
    [
      {"some":"json"}
    ]
    """

  Scenario: Bad request
    When I request REST endpoint with method "DELETE" and path "/bad-request"

    Then I should have a bad request response

    And I should have a response with following JSON body
    """
    {"error":"oops"}
    """

  Scenario: POST with body
    When I request REST endpoint with method "POST" and path "/with-body" and body
    """
    [
      {"some":"json"}
    ]
    """

    Then I should have an OK response

    And I should have a response with following JSON body
    """
    {"status":"ok"}
    """

  Scenario: Successful DELETE Request with no content
    When I request REST endpoint with method "DELETE" and path "/delete-something"

    Then I should have a no content response

  Scenario: Successful DELETE Request with code 204
    When I request REST endpoint with method "DELETE" and path "/delete-something"

    Then I should have a response with following code "204"
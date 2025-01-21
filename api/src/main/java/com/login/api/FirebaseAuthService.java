package com.login.api;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.util.HashMap;
import java.util.Map;
@Service
public class FirebaseAuthService {
    private static final String FIREBASE_AUTH_URL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=()";
    private final String apiKey = "()"; // Replace with your Firebase Web API Key

    public String loginWithEmailAndPassword(String email, String password) {
        RestTemplate restTemplate = new RestTemplate();

        // Request body
        Map<String, Object> requestBody = new HashMap<>();
        requestBody.put("email", email);
        requestBody.put("password", password);
        requestBody.put("returnSecureToken", true);

        HttpHeaders headers = new HttpHeaders();
        headers.set("Content-Type", "application/json");

        HttpEntity<Map<String, Object>> entity = new HttpEntity<>(requestBody, headers);

        try {
            ResponseEntity<Map> response = restTemplate.exchange(
                    FIREBASE_AUTH_URL + apiKey,
                    HttpMethod.POST,
                    entity,
                    Map.class
            );

            // Extract the ID token from the response
            Map<String, Object> responseBody = response.getBody();
            assert responseBody != null;
            return (String) responseBody.get("idToken");

        } catch (Exception e) {
            throw new RuntimeException("Failed to login: " + e.getMessage());
        }
    }
    public String registerWithEmailAndPassword(String email, String password) {
        RestTemplate restTemplate = new RestTemplate();

        // Request body
        Map<String, Object> requestBody = new HashMap<>();
        requestBody.put("email", email);
        requestBody.put("password", password);
        requestBody.put("returnSecureToken", true);

        HttpHeaders headers = new HttpHeaders();
        headers.set("Content-Type", "application/json");

        HttpEntity<Map<String, Object>> entity = new HttpEntity<>(requestBody, headers);

        try {
            ResponseEntity<Map> response = restTemplate.exchange(
                    "https://identitytoolkit.googleapis.com/v1/accounts:signUp?key="
            + apiKey,
                    HttpMethod.POST,
                    entity,
                    Map.class
            );

            Map<String, Object> responseBody = response.getBody();
            assert responseBody != null;
            return (String) responseBody.get("idToken");

        } catch (Exception e) {
            throw new RuntimeException("Failed to register: " + e.getMessage());
        }
    }
}


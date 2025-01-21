package com.login.api;

import com.google.firebase.auth.FirebaseAuth;
import com.google.firebase.auth.FirebaseAuthException;
import com.google.firebase.auth.UserRecord;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class logincontroller {

    @Autowired
    private FirebaseAuthService firebaseAuthService;

    @PostMapping("/login")
    public ResponseEntity<String> login(@RequestParam String email, @RequestParam String password) {
        try {
            String idToken = firebaseAuthService.loginWithEmailAndPassword(email, password);
            return ResponseEntity.ok("Login successful. ID Token: " + idToken);
        } catch (Exception e) {
            return ResponseEntity.badRequest().body("Login failed: " + e.getMessage());
        }
    }
    @PostMapping("/register")
    public ResponseEntity<String> register(@RequestParam String email, @RequestParam String password) {
        try {
            String idToken = firebaseAuthService.registerWithEmailAndPassword(email, password);
            return ResponseEntity.ok("Registration successful. ID Token: " + idToken);
        } catch (Exception e) {
            return ResponseEntity.badRequest().body("Registration failed: " + e.getMessage());
        }
    }
    }


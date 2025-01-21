package com.login.api;

import com.google.auth.oauth2.GoogleCredentials;
import com.google.firebase.FirebaseApp;
import com.google.firebase.FirebaseOptions;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;

@SpringBootApplication
public class ApiApplication {

	public static void main(String[] args) throws IOException {
		SpringApplication.run(ApiApplication.class, args);
		FileInputStream serviceAccount =
				new FileInputStream("C:\\Users\\Admin\\Desktop\\ruman\\api\\src\\main\\resources\\serviceAccountKey.json");

		FirebaseOptions options = new FirebaseOptions.Builder()
				.setCredentials(GoogleCredentials.fromStream(serviceAccount))
				.build();
		FirebaseApp.initializeApp(options);
	}

}

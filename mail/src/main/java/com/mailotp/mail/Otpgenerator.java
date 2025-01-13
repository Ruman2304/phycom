package com.mailotp.mail;

import java.util.Random;

public class Otpgenerator {
    public static String generateOtp(int length) {
        Random random = new Random();
        StringBuilder otp = new StringBuilder();
        for (int i = 0; i < length; i++) {
            otp.append(random.nextInt(10)); // Generate a digit between 0-9
        }
        return otp.toString();
    }

}

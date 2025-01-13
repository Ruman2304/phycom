package com.mailotp.mail;

import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class Otpstore {
    private static final ConcurrentHashMap<String, String> otpStore = new ConcurrentHashMap<>();
    private static final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);

    public static void storeOtp(String email, String otp, int expirySeconds) {
        otpStore.put(email, otp);
        scheduler.schedule(() -> otpStore.remove(email), expirySeconds, TimeUnit.SECONDS);
    }

    public static boolean verifyOtp(String email, String otp) {
        return otp.equals(otpStore.get(email));
    }
}

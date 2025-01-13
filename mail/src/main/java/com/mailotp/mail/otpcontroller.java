package com.mailotp.mail;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
@RestController
@RequestMapping("api/v1/otp")
public class otpcontroller {
         @Autowired
        private Emailservice emailservice;

        @PostMapping("/send")
        public String sendotp(@RequestParam String email){
        String otp = Otpgenerator.generateOtp(6); // Generate a 6-digit OTP
        Otpstore.storeOtp(email, otp, 300); // Store OTP for 5 minutes
        emailservice.sendOtp(email, "Your OTP Code", otp);
        return "OTP sent successfully to " + email;
    }
    @PostMapping("/verify")
    public String verifyOtp(@RequestParam String email, @RequestParam String otp) {
        if (Otpstore.verifyOtp(email, otp)) {
            return "OTP verified successfully!";
        } else {
            return "Invalid or expired OTP.";
        }
    }

}

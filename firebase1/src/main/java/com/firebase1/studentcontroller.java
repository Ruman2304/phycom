package com.firebase1;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.concurrent.ExecutionException;

@RestController
public class studentcontroller {
    public studentcontroller(com.firebase1.studentservice studentservice) {
        this.studentservice = studentservice;
    }
    private studentservice studentservice;

    @PostMapping("/create")
    public String createcrud(@RequestBody student student) throws ExecutionException, InterruptedException {
        return studentservice.createCRUD(student);
    }
    @PutMapping("/update")
    public String updatecrud(@RequestBody student student){
        return studentservice.updateCRUD(student);
    }
    @DeleteMapping ("/delete")
    public String deletecrud(@RequestParam String documentid){
        return studentservice.deleteCRUD(documentid);
    }
    @GetMapping("/get")
    public student getcrud(@RequestParam String documentid) throws ExecutionException, InterruptedException {
        return studentservice.getCRUD(documentid);
    }
    @GetMapping("/test")
    public ResponseEntity<String> testGetEndpoint()
    { return ResponseEntity.ok("Test Get Endpoint is working");
    }
//    @GetMapping("/get")
//    public student getstud(@RequestParam String documentid) throws ExecutionException, InterruptedException {
//        return  studentservice.getcrud(documentid);
//    }
//    @PostMapping("/create")
//    public String createcrud(@RequestBody student student) throws ExecutionException, InterruptedException {
//        return studentservice.createcrud(student);
//    }
//    @PutMapping("/update")
//    public String updatestud(@RequestParam student student) throws ExecutionException, InterruptedException {
//        return  studentservice.updatestudent(student);
//    }
//    @DeleteMapping ("/delete")
//    public String updatestud(@RequestParam String documentid) throws ExecutionException, InterruptedException {
//        return studentservice.deletetudent(documentid);
//    }
}

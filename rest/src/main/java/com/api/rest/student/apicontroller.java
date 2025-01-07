package com.api.rest.student;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import java.util.List;

@RestController
public class apicontroller {
    @Autowired
    private studentrepository studentrepository;

    @GetMapping(path = "/students")
    public List<Student> getstudent(){
        return studentrepository.findAll();
    }
    @PostMapping("/save")
    public String studentsave(@RequestBody Student student){
        studentrepository.save(student);
        return "saved";
    }

    @GetMapping("/")
    public String hello(){
   return "welcome";}

    @PutMapping("/update/{id}")
    public String updatestud(@PathVariable Long id,@RequestBody Student S){
        Student updatestudent= studentrepository.findById(id).get();
        updatestudent.setName(S.getName());
        updatestudent.setAddress(S.getAddress());
        updatestudent.setAge(S.getAge());
        studentrepository.save(updatestudent);
        return "update success";
    }
    @DeleteMapping("delete/{id}")
    public String deleteuser(@PathVariable Long id){
        Student deletestudent=studentrepository.findById(id).get();
        studentrepository.delete(deletestudent);
        return  "student deleted of id"+id;
    }

}

package com.api.rest.student;

import org.springframework.data.jpa.repository.JpaRepository;

public interface studentrepository extends JpaRepository<Student,Long> {
}

package com.sapelkinav.daphelloworld;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class DaphelloworldApplication {

    public static void main(String[] args) {
        var helloWorld = "Hello World!";
        System.out.println(helloWorld);
        SpringApplication.run(DaphelloworldApplication.class, args);
    }

}

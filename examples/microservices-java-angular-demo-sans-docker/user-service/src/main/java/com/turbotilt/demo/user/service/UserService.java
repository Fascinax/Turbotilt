package com.turbotilt.demo.user.service;

import com.turbotilt.demo.user.model.User;
import com.turbotilt.demo.user.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class UserService {

    private final UserRepository userRepository;
    private final RabbitTemplate rabbitTemplate;
    
    private static final String USER_EXCHANGE = "user.exchange";
    private static final String USER_CREATED_KEY = "user.created";
    private static final String USER_UPDATED_KEY = "user.updated";
    
    public List<User> findAllUsers() {
        return userRepository.findAll();
    }
    
    public Optional<User> findUserById(Long id) {
        return userRepository.findById(id);
    }
    
    public User saveUser(User user) {
        boolean isNew = user.getId() == null;
        User savedUser = userRepository.save(user);
        
        // Envoyer un message sur RabbitMQ
        if (isNew) {
            rabbitTemplate.convertAndSend(USER_EXCHANGE, USER_CREATED_KEY, savedUser);
        } else {
            rabbitTemplate.convertAndSend(USER_EXCHANGE, USER_UPDATED_KEY, savedUser);
        }
        
        return savedUser;
    }
    
    public void deleteUser(Long id) {
        userRepository.deleteById(id);
    }
}

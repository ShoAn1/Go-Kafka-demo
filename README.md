# Go-Kafka-demo
Steps to deploy 

ON DOCKER

DOWNLOAD AND INSTALL DOCKER AND DOCKER-COMPOSE 

      >Go to /application

      >Run Command
      -docker-compose --build
      -this will create application image

      >Run command
      -docker-compose up
      -to start the application 
      

ON Minicube 

DOWNLOAD AND SETUP MINIKUBE 

->1.Go to deployments

->2. Run below command
    
    >kubectl apply -f zookeper-deployment.yml
    
    >kubectl apply -f zookeeper-service.yml
    
    >kubectl apply -f kafka-claim0-persistentvolumeclaim.yaml
     
    >kubectl apply -f kafka-service.yaml
  
    >kubectl apply -f kafka-deployment.yaml
    
    >kubectl apply -f application1-claim0-persistentvolumeclaim.yaml
    
    >kubectl apply -f application1-service.yaml
    
    >kubectl apply -f application1-deployment.yaml

#
MiniCloud it 's a infrastructure as container (IAC) App for provisionnig and updating and watching and droping your infra it use terraform for setup a cluster of master and worker node and consul for container microservices discoverry and registration and it can accept ansible playbook path to run with in each container in infra,client should interact with miniCloudCore using CLI and a config file.
Dependencies:GRPC,golang dockerAPI client,consul API client,ansible go client.
#
Install:
git clone https://github.com/Larouimohammed/miniCloud.git
#
make grpc // generate proto code 

#
make dockerinstall // install docker 
#
make ansibleinstall // install ansible with in miniCloud Server
#
make consulinstall // run consul server and agent in docker
#
make run //run miniCloudCore Server
#
make apply or make update or make drop or make watch //run client

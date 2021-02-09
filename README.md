# docker-sdk
Resulta que el término clave a buscar no era tanto ‘docker in docker’, sino ‘docker api and sdks’. Si bien vulgarmente la gente habla de docker in docker, no existe esa terminología en la documentación oficial de docker. Lo que sí existe en la documentación es la API que puede ser accedida como cualquier API mediante http. Adicionalmente, se presentan SDKs para Go y Python.

Dicho esto, encontré artículos que explican cómo correr ‘docker in docker’, y ahí es donde me tranquilicé, porque básicamente se manejan 3 opciones. [El artículo al pie de página](https://devopscube.com/run-docker-in-docker/) es un buen resumen. El problema con esto es que no veo la forma de hacer correr esto desde un software. O sea, estos métodos necesitan una terminal, a lo cual no vamos tener acceso desde nuestros nodos porque queremos que sea automatizado desde un nodo especializado en esta tarea. Y si bien no es para nada imposible, tampoco es del todo elegante que dicho nodo ejecute comando de CLI…

Las referencias a seguir son:
- https://docs.docker.com/engine/api/sdk/examples/
- https://pkg.go.dev/github.com/docker/docker/client 

Después de haber probado estas cosas, me doy cuenta de que el [primer artículo que referencié y que dije que no sirve, en realidad sirve](https://devopscube.com/run-docker-in-docker/). Detallo la experiencia a continuación.
1. Dockerfile. Usé el [Dockerfile presente en el Go de Docker hub](https://hub.docker.com/_/golang). Antes de esto probé otros de unos tutoriales, pero no se llevaron bien con la importación de módulos de GitHub.
2. Seguí el [ejemplo oficial más básico de Docker Go SDK](https://pkg.go.dev/github.com/docker/docker/client). El problema surgió cuando quise ejecutar esto en Docker. Hice build del container. Todo OK. Pero cuando quise hacer run me topé con el error:
```
panic: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
goroutine 1 [running]:
main.main()
	/go/src/app/main.go:19 +0x2b9
```
Acá es donde entra en juego [el primer artículo que mencioné](https://devopscube.com/run-docker-in-docker/): siguiendo el primero de los 3 métodos presentados, le agregué a docker run el flag `-v /var/run/docker.sock:/var/run/docker.sock`. Con ese flag ya no tuve ningún error, pero tampoco input, lo que era esperable porque no había otros contenedores corriendo (docker ps no listaba ningún contenedor). Luego de dejar corriendo una terminal de ubuntu desde docker, volví a correr el container de Go, y listó el container de ubuntu según lo esperado.
3. El próximo paso es levantar algún container desde este programita...

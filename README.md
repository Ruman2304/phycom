Lax-Alerts Appengine Layer
=============================================
For additional details see google quick start and related documentation [https://cloud.google.com/endpoints/docs/frameworks/java/quickstart-frameworks-java].

# Deployment Steps,

1. `mvn clean package`
2. `mvn exec:java -DGetSwaggerDoc`
3. `gcloud endpoints services deploy openapi.json`
4. `mvn appengine:deploy`

## Stage Environment
The pom.xml file contains stage specic settings overrides and a BuildTimeConstants file is dynamically build at compile time. With this in mind deploying to stage is as simple as preparing a clean build and using appengine:deploy with the stage profile enabled.  The deployment project, etc. is set within the pom for you.
variant
```
mvn clean package -P stage
mvn appengine:deploy -P stage

`mvn exec:java -DGetSwaggerDoc -P stage`
`gcloud endpoints services deploy openapi.json --project lax-gateway-stage`

```


# Old method
`mvn clean package && mvn exec:java -DGetSwaggerDoc && gcloud service-management deploy openapi.json`

# New method - does not yet work correctly
`mvn clean package && mvn clean package endpoints-framework:openApiDocs -DskipTests && gcloud endpoints services deploy target/openapi-docs/openapi.json`

# Developer Sandbox

## New Plugin
Run `mvn appengine:run` to start local sandbox.


# Update Composite Indexes
`gcloud beta emulators datastore start --data-dir ./src/main/webapp/WEB-INF/`

# Push Composite Indexes
Run:
`gcloud datastore indexes create src/main/webapp/WEB-INF/index.yaml  --project lax-gateway`
`gcloud datastore indexes create src/main/webapp/WEB-INF/index.yaml  --project lax-gateway-stage`


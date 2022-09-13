# Dynamic Load Balancer configuration for Cloud Run services
When adding additional Cloud Run services development teams enjoy focusing on application code rather than infrastructure.  In some cases a load balancer may be requred.  You may desire to use path based routing or the ablility to protect the Cloud Run application APIs using Apigee. 

This is example of how to use the `URL masking` feature of the internal load balancer to call cloud run services dynamically without additional load balancer  configuration for additional cloud run services.

For the most part this is already [documented](https://cloud.google.com/load-balancing/docs/l7-internal/setting-up-l7-internal-serverless). This example incorporates the [URL mask feature](https://cloud.google.com/load-balancing/docs/l7-internal/setting-up-l7-internal-serverless#using-url-mask) allowing us to understand how the URL path mapping will work.

Using this solution will require the name of the service to be included in your URL path, for example imagine you have an employee service and it a has a find function mapped to the url `/find`.  After implementing this solution the URL will be `/employee/find`.  This will require an update to the path prefix in the application to include the service name.

## Demo
* Follow the instructions in [Set up an internal HTTP(S) load balancer with Cloud Run](https://cloud.google.com/load-balancing/docs/l7-internal/setting-up-l7-internal-serverless) with the following exceptions bulleted below.  Use the `gcloud` version of the instructions and HTTP protocol. The instructions can be adjusted for external the external loadbalancer and HTTPS.

* We will use a dynamic go lang application to help understand this. Instead of using the service defined at [Deploy a Cloud Run service](https://cloud.google.com/load-balancing/docs/l7-internal/setting-up-l7-internal-serverless#deploy_serverless_app) use the services defined below.
```
gcloud run deploy employee \
  --source=. \
  --allow-unauthenticated \
  --ingress=all \
  --region=us-central1 

gcloud run deploy product \
  --source=. \
  --allow-unauthenticated \
  --ingress=all \
  --region=us-central1 
```


* In [Create the load balancer](https://cloud.google.com/load-balancing/docs/l7-internal/setting-up-l7-internal-serverless#creating_the_load_balancer) section use these instructions to create the URL mask instead of is defined in the first setp (`gcloud compute network-endpoint-groups...` )
```
gcloud compute network-endpoint-groups create SERVERLESS_NEG_NAME \
    --region=REGION \
    --network-endpoint-type=serverless \
    --cloud-run-url-mask="/<service>"
```

* In [Send traffic to the load balancer](https://cloud.google.com/load-balancing/docs/l7-internal/setting-up-l7-internal-serverless#send_traffic_to_the_load_balancer) use the following URLs to understand how the mapping works
```
curl 10.1.2.99/employee  # results in the `root` function being called in the employee service
curl 10.1.2.99/employee/find  # results in `serviceFind` function being called in the employee service
curl 10.1.2.99/product  # results in the `root` function being called in the product service
curl 10.1.2.99/product/find  # results in `serviceFind` function being called in the product service
curl 10.1.2.99/find  # results in an empty reply - no cloud run find service defined
```


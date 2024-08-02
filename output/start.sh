# RegService
#echo Starting first RegService
/workspaces/RLAD/output/large-scale-workshop /workspaces/RLAD/services/reg-service/service/RegServiceRoot.yaml
#sleep 1

#echo Starting second
/workspaces/RLAD/output/large-scale-workshop /workspaces/RLAD/services/reg-service/service/AnotherRegService.yaml


# CacheService
/workspaces/RLAD/output/large-scale-workshop /workspaces/RLAD/services/cache-service/service/CacheServiceRoot.yaml
/workspaces/RLAD/output/large-scale-workshop /workspaces/RLAD/services/cache-service/service/AnotherCacheService.yaml

# TestService
/workspaces/RLAD/output/large-scale-workshop /workspaces/RLAD/services/test-service/service/TestService.yaml

# ...
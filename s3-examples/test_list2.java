package com.deeproute.dataocean.file.galaxy;

import com.deeproute.dataocean.bootstarter.asynctask.constant.AsyncTaskConfigStatus;
import com.deeproute.dataocean.bootstarter.mongodb.sharding.ShardingMongoTemplate;
import com.deeproute.dataocean.bootstarter.mongodb.sharding.util.ShardingUtils;
import com.deeproute.dataocean.bootstarter.s3.S3ClientFactory;
import com.deeproute.dataocean.bootstarter.s3.model.S3Config;
import com.deeproute.dataocean.common.App;
import com.deeproute.dataocean.common.annotation.AsyncCall;
import com.deeproute.dataocean.common.model.ObjectHolder;
import com.deeproute.dataocean.common.model.SelfBean;
import com.deeproute.dataocean.common.river.node.Rivers;
import com.deeproute.dataocean.common.river.plugin.scheduler.Schedulers;
import com.deeproute.dataocean.common.utils.CollectionUtil;
import com.deeproute.dataocean.common.utils.JsonUtils;
import com.deeproute.dataocean.common.utils.RetryUtil;
import com.deeproute.dataocean.file.galaxy.api.module.filestorage.model.bo.storageimpl.S3StorageDetail;
import com.deeproute.dataocean.file.galaxy.constant.MongoConst;
import com.deeproute.dataocean.file.galaxy.module.drfilesharding.infrastructure.cache.DrFileShardingCache;
import com.deeproute.dataocean.file.galaxy.module.file.constant.DrFileConst;
import com.deeproute.dataocean.file.galaxy.module.file.model.bo.StorageInfo;
import com.deeproute.dataocean.file.galaxy.module.file.model.entity.BaseDrFile;
import com.deeproute.dataocean.file.galaxy.module.file.model.entity.DrFile;
import com.deeproute.dataocean.file.galaxy.module.file.service.IDrFileQueryService;
import com.deeproute.dataocean.file.galaxy.module.filledelete.model.entity.TrashDrFile;
import com.deeproute.dataocean.file.galaxy.module.versioneddrfile.model.entity.VersionedDrFile;
import com.mongodb.bulk.BulkWriteResult;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.experimental.SuperBuilder;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.codec.digest.DigestUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.mongodb.core.BulkOperations;
import org.springframework.data.mongodb.core.query.Criteria;
import org.springframework.data.mongodb.core.query.Query;
import org.springframework.data.mongodb.core.query.Update;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.data.util.Pair;
import org.springframework.stereotype.Service;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.ListObjectsRequest;
import software.amazon.awssdk.services.s3.model.ListObjectsResponse;
import software.amazon.awssdk.services.s3.model.ListObjectsV2Request;
import software.amazon.awssdk.services.s3.model.ListObjectsV2Response;
import software.amazon.awssdk.services.s3.model.S3Object;

import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.Set;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;

/**
 * @author Jianhui Jiang
 * @description
 * @email jianhuijiang@deeproute.ai
 * @date 2025/3/20 11:21
 */
@Slf4j
@Service
public class Ssd11TempService extends SelfBean<Ssd11TempService> {
    //    ssd12: 10.3.10.158-172 （15）10.3.10.139-141 （3）
//    ak: 090OVX2XJ9P709DOCM7F
//    sk: 8nWw976gATZFQz2uMkwoZxdKgTh1W4ifiHjmAZz6
    public static final String[] endpoints = {
            "http://10.3.10.139:80",
            "http://10.3.10.140:80",
            "http://10.3.10.141:80",

            "http://10.3.10.158:80",
            "http://10.3.10.159:80",
            "http://10.3.10.160:80",
            "http://10.3.10.161:80",
            "http://10.3.10.162:80",
            "http://10.3.10.163:80",
            "http://10.3.10.164:80",
            "http://10.3.10.165:80",
            "http://10.3.10.166:80",
            "http://10.3.10.167:80",
            "http://10.3.10.168:80",
            "http://10.3.10.169:80",
            "http://10.3.10.170:80",
            "http://10.3.10.171:80",
            "http://10.3.10.172:80",
    };
    public static final String SSD12_AK = "090OVX2XJ9P709DOCM7F";
    public static final String SSD12_SK = "8nWw976gATZFQz2uMkwoZxdKgTh1W4ifiHjmAZz6";

    public static final String SSD11_AK = "SVVPV58E6WQ5EE6NY2JB";
    public static final String SSD11_SK = "4FmBvshd2BRNP6UuGJpvnFYGaNrAnEbTiXLkfQZZ";

    @Autowired
    private S3ClientFactory s3ClientFactory;
    @Autowired
    private IDrFileQueryService drFileQueryService;
    @Autowired
    private DrFileShardingCache shardingCache;
    @Autowired
    private ShardingMongoTemplate template;
    @Autowired
    private StringRedisTemplate redisTemplate;

    private S3Client client;

    public static void main(String[] args) {

        S3ClientFactory factory = S3ClientFactory.create();
        long l = System.currentTimeMillis();

        try {
            S3Config config = S3Config.builder()
                    .endpoint("http://10.3.10.145:80")
                    .accessKey(SSD11_AK)
                    .secretKey(SSD11_SK)
                    .readTimeout(TimeUnit.SECONDS.toMillis(10))
                    .build();
            S3Client client1 = factory.getS3client(config);


            ListObjectsV2Response res = client1.listObjectsV2(ListObjectsV2Request.builder()
                    .bucket("prod-trip-1")
                    .maxKeys(1000)
                    .continuationToken("9fffff0a-d028-4eb2-ab8f-a5095f03db32")
                    .build());
            for (S3Object obj : res.contents()) {
                System.out.println(obj.key());
            }
        } catch (Throwable e) {
            System.out.println(TimeUnit.MILLISECONDS.toSeconds(System.currentTimeMillis() - l));
            throw e;
        }

    }

    private S3Client getClient() {
        if (client != null) {
            return client;
        }

        S3Config config = S3Config.builder()
                .endpoint("http://10.3.10.143:80")
                .accessKey(SSD11_AK)
                .secretKey(SSD11_SK)
                .build();
        client = s3ClientFactory.getS3client(config);

        return client;
    }

    public void changeSsd11(String bucket, String marker, int maxKey) {
        AtomicInteger remain = new AtomicInteger(maxKey);
        AtomicLong scannedCount = new AtomicLong(0l);
        ObjectHolder<String> theMarker = ObjectHolder.of(marker);
        try {
            while (true) {
                log.info("change ssd11, scan for {}, marker: {}, maxKey: {}, scanned: {}", bucket, theMarker.get(), maxKey, scannedCount.get());
                boolean result = RetryUtil.retryGet(10, times -> {
                    Integer batch = App.getProperty("tmp.ssd11.batch", Integer.class, 1000);
                    if (maxKey > 0) {
                        batch = Math.min(remain.intValue(), batch);
                    }
                    if (batch <= 0) {
                        return false;
                    }
                    ListObjectsV2Response response = getClient().listObjectsV2(ListObjectsV2Request.builder()
                            .bucket(bucket)
                            .maxKeys(batch)
                            .continuationToken(theMarker.get())
                            .build()
                    );

                    if (response.hasContents()) {
                        remain.addAndGet(-response.contents().size());
                        scannedCount.addAndGet(response.contents().size());
                        System.out.println(response.contents().size());
                        Rivers.fromList(response.contents())
                                .batch(10)
                                .scheduleOn(Schedulers.io(10))
                                .foreach(keys -> {
                                    Set<String> dataKeys = Rivers.fromList(keys)
                                            .map(obj -> {
                                                String key = obj.key();
                                                String oldDataKey = String.join(":", key, bucket, SSD11_AK);
                                                String newDataKey = DigestUtils.md2Hex(oldDataKey);

                                                if (log.isDebugEnabled()) {
                                                    log.debug("change ssd11, key: {}", key);
                                                    log.debug("change ssd11, newDataKey: {}", newDataKey);
                                                    log.debug("change ssd11, oldDataKey: {}", oldDataKey);
                                                }
                                                return new String[]{oldDataKey, newDataKey};
                                            })
                                            .flatmap(Rivers::fromArray)
                                            .toSet();
                                    getSelf().changeByDataKeys(Bo.builder().bucket(bucket).dataKeys(dataKeys).build());
                                });
                    }
                    if (!response.isTruncated() || response.nextContinuationToken() == null) {
                        return false;
                    }

                    theMarker.set(response.nextContinuationToken());
                    return true;
                }, Throwable.class);
                if (!result) {
                    break;
                }
            }

            log.info("change ssd11, scan finished for {}", bucket);
        } catch (Throwable e) {
            log.warn("change ssd11, scan failed，bucket: {}, marker:{}", bucket, theMarker.get(), e);
        }
    }

    public void changeByDataKeysSync(Set<String> dataKeys) {
        changeByDataKeys(Bo.builder().dataKeys(dataKeys).build());
    }

    @Data
    @AllArgsConstructor
    @NoArgsConstructor
    @SuperBuilder
    public static class Bo {
        private String bucket;
        private Set<String> dataKeys;
    }

    @AsyncCall(
            taskId = @AsyncCall.TaskId(generateClass = AsyncCall.UUIDTaskIdGenerator.class),
            status = AsyncTaskConfigStatus.OFF,
            retryTimes = -1,
            completeTtl = @AsyncCall.Duration,
            cron = "0/1 * * * * ?"
    )
    protected void changeByDataKeys(Bo bo) {
        Set<String> dataKeys = bo.dataKeys;
        List<BaseDrFile> drFiles = drFileQueryService.findAllDrFilesByDataKey(dataKeys);
        recordNotExistDataKey(drFiles, dataKeys);
        if (log.isDebugEnabled()) {
            log.debug("change ssd11, dataKeys: {}", JsonUtils.serialize(dataKeys));
            log.debug("change ssd11, drFiles size: {}", drFiles.size());
            for (BaseDrFile drFile : drFiles) {
                log.debug("change ssd11, {}: {}", drFile.getClass().getSimpleName(), JsonUtils.serialize(drFile));
            }
        }
        if (drFiles.isEmpty()) {
            return;
        }
        Map<Object, List<Pair<Query, Update>>> groups = CollectionUtil.group(drFiles, this::group, drFile -> {
            Random random = new Random();
            for (StorageInfo info : drFile.getStorageInfo()) {
                if (dataKeys.contains(info.getDetail().getDataKey())) {
                    S3StorageDetail detail = (S3StorageDetail) info.getDetail();
                    detail.setEndpoint(endpoints[random.nextInt(endpoints.length)]);
                    detail.setAccessKey(SSD12_AK);
                    detail.setSecretKey(SSD12_SK);

                    info.setDataKey(detail.getDataKey());
                    info.setCluster("bigdata-ssd-12");
                    info.setName("BIGDATA_SSD_12");
                }
            }

            Query query = new Query();
            query.addCriteria(Criteria.where(MongoConst.FIELD_ID).is(drFile.getId()));

            Update update = new Update();
            update.set(BaseDrFile.Fields.storageInfo, drFile.getStorageInfo());
            return Pair.of(query, update);
        });

        if (log.isDebugEnabled()) {
            for (Object o : groups.keySet()) {
                log.debug("change ssd11, group key: {}", o);
            }
        }

        groups.forEach((k, group) -> {
            Class entityClass = k instanceof String ? DrFile.class : (Class) k;
            try {
                if (entityClass == DrFile.class) {
                    ShardingUtils.mark(DrFileConst.COLLECTION_NAME_DR_FILE, (String) k);
                }
                BulkWriteResult result = template.bulkOps(BulkOperations.BulkMode.UNORDERED, entityClass).updateOne(group).execute();
                if (result.getModifiedCount() != group.size()) {
                    throw new RuntimeException("更新条数不一致！");
                }
            } finally {
                ShardingUtils.clearMark(DrFileConst.COLLECTION_NAME_DR_FILE);
            }
        });
    }


    private void recordNotExistDataKey(List<BaseDrFile> drFiles, Set<String> dataKeys) {
        HashSet<String> notExist = new HashSet<>();
        for (String dataKey : dataKeys) {
            if (dataKey.contains(":")) {
                notExist.add(dataKey.split(":")[0]);
            }
        }

        for (BaseDrFile drFile : drFiles) {
            for (StorageInfo storageInfo : drFile.getStorageInfo()) {
                if (!(storageInfo.getDetail() instanceof S3StorageDetail)) {
                    continue;
                }
                S3StorageDetail detail = (S3StorageDetail) storageInfo.getDetail();
                if (!SSD11_AK.equals(detail.getAccessKey())) {
                    continue;
                }
                notExist.remove(detail.getKey());
            }
        }

        if (!notExist.isEmpty()) {
            redisTemplate.opsForSet().add("dataocean-file-galaxy:ssd11:no-drfile-key", notExist.toArray(new String[0]));
        }
    }

    private Object group(BaseDrFile drFile) {
        if (drFile instanceof TrashDrFile) {
            return TrashDrFile.class;
        }
        if (drFile instanceof VersionedDrFile) {
            return VersionedDrFile.class;
        }
        String shardingId = shardingCache.getShardingId(drFile.getRealTreeId(), null);
        return shardingId;
    }
}

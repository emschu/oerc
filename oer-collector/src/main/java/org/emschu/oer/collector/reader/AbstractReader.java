package org.emschu.oer.collector.reader;

/*-
 * #%L
 * oer-server
 * %%
 * Copyright (C) 2019 emschu[aet]mailbox.org
 * %%
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * #L%
 */

import org.emschu.oer.collector.reader.parser.*;
import org.emschu.oer.collector.service.ChannelService;
import org.emschu.oer.collector.util.DateRangeUtil;
import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.model.TvShow;
import org.emschu.oer.core.model.UpdateControlInterface;
import org.emschu.oer.core.model.repository.ProgramEntryRepository;
import org.emschu.oer.core.model.repository.TvShowRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;
import org.springframework.stereotype.Service;

import javax.annotation.PreDestroy;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeParseException;
import java.time.temporal.ChronoUnit;
import java.util.*;
import java.util.concurrent.LinkedBlockingDeque;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;
import java.util.logging.Level;
import java.util.logging.Logger;

/**
 * abstract class representing a reader. encapsulates common logic of all readers.
 * manages event handling
 */
@Configuration
@Service
@PropertySource("classpath:oer.properties")
public abstract class AbstractReader implements ApplicationContextAware {

    private final Logger LOG = Logger.getLogger(AbstractReader.this.getClass().getName());
    private static final int DAY_BLOCK_LIST_SIZE = 25;
    private boolean isInitialized = false;
    private boolean isRunning = false;
    private List<CustomParser> customParsers = new ArrayList<>();

    @Autowired
    private ApplicationContext context;

    @Autowired
    private ChannelService channelService;

    // repos
    @Autowired
    private ProgramEntryRepository programEntryRepository;

    @Autowired
    private TvShowRepository tvShowRepository;
    private List<ThreadPoolExecutor> threadPoolExecutors = new ArrayList<>();

    // properties
    @Value(value = "${oer.collector.update_mode}")
    private boolean updateMode;

    @Value(value = "${oer.collector.update_mode_force}")
    private boolean updateModeForce;

    @Value(value = "${oer.collector.enable_tv_show_collect}")
    private boolean collectTvShows;

    @Value(value = "${oer.collector.enable_program_entry_collect}")
    private boolean collectProgramEntries;

    @Value(value = "${oer.collector.collect_future_program_days_max}")
    private String collectDaysInFutureMax;

    @Value(value = "${oer.collector.collect_past_program_days_max}")
    private String collectDaysInPastMax;

    @Value(value = "${debug}")
    private boolean isDebug;

    @Value(value = "${oer.collector.invalidate_update_hours}")
    private String invalidationHours;

    @Value(value = "${oer.collector.skip_ard}")
    private boolean skipArd;

    @Value(value = "${oer.collector.skip_zdf}")
    private boolean skipZdf;

    @Value(value = "${oer.collector.skip_orf}")
    private boolean skipOrf;

    @Value(value = "${oer.collector.mass_mode}")
    private boolean isMassMode;

    @Value(value = "${oer.collector.start_date}")
    private String startDate;

    @Value(value = "${oer.collector.end_date}")
    private String endDate;

    @Value(value = "${oer.collector.core_thread_pool_size}")
    private String coreThreadPoolSizeProperty;

    @Value(value = "${oer.collector.max_thread_pool_size}")
    private String maxThreadPoolSizeProperty;
    private int maxThreadPoolSize;
    private int coreThreadPoolSize;

    public void checkConfiguration() {
        if (isInitialized) {
            // avoid duplicate execution of subclass postconstruct calls
            return;
        }
        LOG.info("Loaded configuration of collector:");
        LOG.info(String.format("oer.collector.update_mode: %s", updateMode));
        LOG.info(String.format("oer.collector.update_mode_force: %s", updateModeForce));
        LOG.info(String.format("oer.collector.enable_tv_show_collect: %s", collectTvShows));
        LOG.info(String.format("oer.collector.enable_program_entry_collect: %s", collectProgramEntries));
        LOG.info(String.format("oer.collector.collect_future_program_days_max: %s", collectDaysInFutureMax));
        LOG.info(String.format("oer.collector.collect_past_program_days_max: %s", collectDaysInPastMax));
        LOG.info(String.format("oer.collector.invalidate_update_hours: %s", invalidationHours));
        LOG.info(String.format("oer.collector.skip_ard: %s", skipArd));
        LOG.info(String.format("oer.collector.skip_zdf: %s", skipZdf));
        LOG.info(String.format("oer.collector.skip_orf: %s", skipOrf));
        LOG.info(String.format("oer.collector.mass_mode: %s", isMassMode));

        try {
            maxThreadPoolSize = Integer.valueOf(this.maxThreadPoolSizeProperty);
            coreThreadPoolSize = Integer.valueOf(this.coreThreadPoolSizeProperty);
        } catch (NumberFormatException nfe) {
            LOG.warning("invalid thread pool config detected:");
            LOG.warning(this.maxThreadPoolSizeProperty);
            LOG.warning(this.coreThreadPoolSizeProperty);
            maxThreadPoolSize = 10;
            coreThreadPoolSize = 5;
        }

        if (getCollectDaysInFutureMax() > 44) {
            throw new IllegalStateException("cannot collect more days in future as 44");
        }
        if (getCollectDaysInPastMax() < 0 || getCollectDaysInFutureMax() < 0) {
            throw new IllegalStateException("invalid negative value for data collection date range");
        }
        if (getInvalidationHours() == 0 || getInvalidationHours() < 0) {
            LOG.warning("Invalidate_update_hours is not set correctly. using fallback '24'");
            invalidationHours = "24";
        }
        if (skipArd && skipZdf && skipOrf) {
            LOG.warning("NOTE: Update for zdf and ard is completely disabled");
        }

        setInitialized();
    }

    private final void setInitialized() {
        isInitialized = true;
    }

    @PreDestroy
    public void destroy() {
        for (ThreadPoolExecutor tpe : threadPoolExecutors) {
            tpe.shutdown();
        }
    }

    // implement these methods in the specific parsers

    /**
     * a program entry parser interface
     *
     * @return a {@link ProgramEntryParserInterface} instance
     */
    public abstract ProgramEntryParserInterface getProgramEntryParser();

    /**
     * a tv show parser interface
     *
     * @return a {@link TvShowParserInterface} instance
     */
    public abstract TvShowParserInterface getTvShowParser();

    /**
     * define an available channel family defined in enum
     *
     * @return an {@link org.emschu.oer.core.model.Channel.AdapterFamily} instance
     */
    public abstract Channel.AdapterFamily getAdapterFamily();

    /**
     * main method of a reader to execute
     *
     * @throws ParserException if anything goes wrong
     * @throws InterruptedException if interrupted
     */
    public void execute() throws ParserException, InterruptedException {
        checkConfiguration();
        if (isRunning) {
            LOG.warning("No execution of " + this.getClass().getName() + ". An update process is already running.");
            return;
        }
        setRunning(true);
        long start = System.currentTimeMillis();
        LOG.info(String.format("Start fetching data of %s", getAdapterFamily()));
        update();
        clear();
        long end = System.currentTimeMillis();
        LOG.info(String.format("End fetching data of %s", getAdapterFamily()));
        LOG.info(String.format("Execution Time: %d ms", (end - start)));
        setRunning(false);
    }

    private void setRunning(boolean isRunning) {
        this.isRunning = isRunning;
    }

    /**
     * update/import process of reader
     */
    private void update() throws ParserException, InterruptedException {
        // check for skipping first
        if (getAdapterFamily().equals(Channel.AdapterFamily.ARD) && skipArd) {
            LOG.warning("ARD update is configured being skipped");
            return;
        }
        if (getAdapterFamily().equals(Channel.AdapterFamily.ZDF) && skipZdf) {
            LOG.warning("ZDF update is configured being skipped");
            return;
        }
        if (getAdapterFamily().equals(Channel.AdapterFamily.ORF) && skipOrf) {
            LOG.warning("ZDF update is configured being skipped");
            return;
        }
        // start with tv shows
        if (collectTvShows) {
            // store tv shows
            LOG.info("Start fetching tv shows...");
            readTvShows();
            LOG.info("End fetching tv shows...");
        } else {
            LOG.info("Collecting tv shows is generally disabled");
        }

        // proceed with program entries
        if (collectProgramEntries) {
            LOG.info("Start fetching program data of channels");
            // store program entries per channel
            List<Channel> channelList = getAllChannels();
            if (channelList == null || channelList.isEmpty()) {
                LOG.warning("null or empty channel list given.");
            } else {
                for (Channel aChannel : channelList) {
                    LOG.info(String.format("Processing channel: '%s'", aChannel.getName()));
                    readProgramEntries(aChannel);
                }
            }
        } else {
            LOG.info("Collecting program entries is generally disabled");
        }

        handleCustomParsers();
    }

    /**
     * internal helper method to run custom parsers
     *
     * @throws ParserException
     */
    private void handleCustomParsers() throws ParserException {
        if (!this.customParsers.isEmpty()) {
            LOG.info(String.format("Detected %d registered custom parsers", this.customParsers.size()));
            final List<LocalDate> datesToFetch = generateDateRangeToFetch();
            for (CustomParser customParser : this.customParsers) {
                // init custom parser
                customParser.setDateRangeList(datesToFetch);
                customParser.setProgramEntryCollectingEnabled(collectProgramEntries);
                customParser.setTvShowCollectingEnabled(collectTvShows);
                customParser.setAdapterFamily(getAdapterFamily());

                LOG.info(String.format("Running: %s", customParser.getClass().getName()));
                customParser.run();
            }
        }
    }

    /**
     * clean up parsers after update
     */
    public void clear() {
        LOG.info("Cleaning reader and its dependent services");
        getTvShowParser().cleanup();
        getProgramEntryParser().cleanup();
        for (CustomParser parser : this.customParsers) {
            parser.cleanup();
        }
        System.gc();
    }

    public List<CustomParser> getCustomParsers() {
        return customParsers;
    }

    /**
     * process a channel
     *
     * @param channel
     * @throws InterruptedException
     */
    private void readProgramEntries(Channel channel) throws InterruptedException {
        List<LocalDate> dayListToFetch = generateDateRangeToFetch();

        // store program entries
        LOG.info("Storing tv program for " + dayListToFetch.size() + " days for channel: " + channel.getName());

        List<List<LocalDate>> dateBlocks = getDateBlocks(dayListToFetch);

        LOG.info("created " + dateBlocks.size() + " date blocks to process");
        final ThreadPoolExecutor channelProgramThreadPool =
                new ThreadPoolExecutor(coreThreadPoolSize, maxThreadPoolSize, 30, TimeUnit.SECONDS, new LinkedBlockingDeque<>(Integer.MAX_VALUE));
        channelProgramThreadPool.prestartAllCoreThreads();

        for (List<LocalDate> dateBlock : dateBlocks) {
            for (LocalDate day : dateBlock) {
                LOG.info("Start fetching day: " + day.format(DateTimeFormatter.ISO_LOCAL_DATE));

                try {
                    handleProgramEntries(channelProgramThreadPool, channel, day);
                } catch (ParserException pe) {
                    LOG.log(Level.WARNING, pe.getMessage(), pe);
                    LOG.throwing(AbstractReader.class.getName(), "readProgramEntries", pe);
                } catch (Exception e) {
                    LOG.log(Level.WARNING, e.getMessage(), e);
                    LOG.throwing(AbstractReader.class.getName(), "readProgramEntries", e);
                    return;
                }

                LOG.info("End fetching day: " + day.format(DateTimeFormatter.ISO_LOCAL_DATE));

                waitMax10Minutes(channelProgramThreadPool);
            }
            LOG.info("Waiting for date block processing");
        }

        channelProgramThreadPool.shutdown();
    }

    /**
     * divide date range day list into a partitioned list with block size DAY_BLOCK_LIST_SIZE
     *
     * @param dayListToFetch
     * @return
     */
    private List<List<LocalDate>> getDateBlocks(List<LocalDate> dayListToFetch) {
        List<List<LocalDate>> dateBlocks = new ArrayList<>();
        for (int i = 0; i <= dayListToFetch.size() / DAY_BLOCK_LIST_SIZE; i++) {
            int upperLimit = (i + 1) * DAY_BLOCK_LIST_SIZE;
            if (upperLimit > dayListToFetch.size()) {
                upperLimit = dayListToFetch.size();
            }
            dateBlocks.add(dayListToFetch.subList(i * DAY_BLOCK_LIST_SIZE, upperLimit));
        }
        return dateBlocks;
    }

    /**
     * method to collect tv shows
     *
     * @throws InterruptedException interrupted
     */
    protected void readTvShows() throws InterruptedException {
        TvShowParserInterface tvShowParser = getTvShowParser();
        if (tvShowParser == null) {
            throw new IllegalStateException("null tv show parser instance given");
        }

        ThreadPoolExecutor tpe = new ThreadPoolExecutor(coreThreadPoolSize, maxThreadPoolSize,
                30, TimeUnit.SECONDS, new LinkedBlockingDeque<>(Integer.MAX_VALUE));
        tpe.prestartAllCoreThreads();

        List<TvShow> tvShowList = detectTvShows();
        List<TvShowPostProcessThread> calls = new ArrayList<>();

        // set end dates as the html does not contain this information
        for (TvShow tvShow : tvShowList) {
            if (tvShow == null) {
                throw new IllegalStateException("null tv show given");
            }
            if (tvShow.getId() == null) {
                calls.add(new TvShowPostProcessThread(tvShow));
            } else {
                // not a new entry!
                if (this.needsUpdate(tvShow)) {
                    // only update data if needed!
                    calls.add(new TvShowPostProcessThread(tvShow));
                }
            }
        }
        if (!calls.isEmpty()) {
            LOG.info("Start fetching of " + calls.size() + " tv show calls");
        }

        for (TvShowPostProcessThread thread : calls) {
            tpe.execute(thread);
        }

        waitMax10Minutes(tpe);
        tpe.shutdown();
    }

    /**
     * waiting loop
     *
     * @param tpe
     * @throws InterruptedException
     */
    private void waitMax10Minutes(ThreadPoolExecutor tpe) throws InterruptedException {
        // 10 min
        final int timeoutSeconds = 600;
        final LocalDateTime executionStartDt = LocalDateTime.now();

        // wait until its finished..
        while (tpe.getActiveCount() > 0) {
            synchronized (this) {
                this.wait(5000);
            }
            if (ChronoUnit.SECONDS.between(executionStartDt, LocalDateTime.now()) > timeoutSeconds) {
                LOG.warning("Timeout reached, stop to wait for tv show threads.");
                // leave loop if timeout is reached
                break;
            }
        }
    }

    /**
     * has at least 1 entry = today
     *
     * @return
     */
    private final List<LocalDate> generateDateRangeToFetch() {
        if (isMassMode) {
            Optional<LocalDate> startDateOptional = getStartDate();
            Optional<LocalDate> endDateOptional = getEndDate();
            if (!startDateOptional.isPresent()) {
                LOG.warning("invalid start date in mass mode given!");
                return new ArrayList<>();
            }
            if (endDateOptional.isPresent()) {
                return DateRangeUtil.dateRangeBetween(startDateOptional.get(), endDateOptional.get());
            }
            return DateRangeUtil.dateRangeBetween(startDateOptional.get(), null);
        }
        return DateRangeUtil.generateDateRangeToFetch(getCollectDaysInFutureMax(), getCollectDaysInPastMax());
    }

    /**
     * collect program entries.
     * NOTE: a developer must
     *
     * @param threadPoolTaskExecutor executor object
     * @param channel channel record
     * @param day the day to fetch
     * @throws ParserException internal error
     */
    protected void handleProgramEntries(ThreadPoolExecutor threadPoolTaskExecutor, Channel channel, LocalDate day) throws ParserException {
        LinkedList<ProgramEntry> linkedTvList = detectProgramEntries(channel, day);
        List<ProgramEntryPostProcessThread> calls = new ArrayList<>();

        // set end dates as the html does not contain this information

        for (ProgramEntry programEntry : linkedTvList) {
            if (programEntry == null) {
                throw new IllegalStateException("null program entry given");
            }
            // zdf program entries are handled differently
            if (programEntry.getAdapterFamily() != Channel.AdapterFamily.ZDF
                    && (programEntry.getStartDateTime() == null || programEntry.getEndDateTime() == null)) {
                LOG.warning("SKIPPING. No start date time or end date time for entry " + programEntry.toString());
                continue;
            }
            if (programEntry.getId() == null) {
                calls.add(new ProgramEntryPostProcessThread(getProgramEntryParser(), programEntryRepository,
                        programEntry, isDebug));
            } else {
                // not a new entry!
                // todo check what happens if .getDescription() is no longer a condition here
                if (programEntry.getDescription() == null || needsUpdate(programEntry)) {
                    // only update data if needed!
                    calls.add(new ProgramEntryPostProcessThread(getProgramEntryParser(), programEntryRepository,
                            programEntry, isDebug));
                }
            }
        }

        if (!calls.isEmpty()) {
            LOG.info("Start fetching of " + calls.size() + " calls");
        }
        for (ProgramEntryPostProcessThread thread : calls) {
            threadPoolTaskExecutor.execute(thread);
        }
    }

    /**
     * returns a filtered list of program entries. only returning records which need to be updated.
     *
     * @param channel
     * @param day
     * @return
     * @throws ParserException
     */
    @SuppressWarnings("unchecked")
    private LinkedList<ProgramEntry> detectProgramEntries(Channel channel, LocalDate day) throws ParserException {
        int entryCounter = 0;
        int errorCounter = 0;
        LinkedList<ProgramEntry> linkedProgramList = new LinkedList<>();

        final ProgramEntryParserInterface programEntryParser = getProgramEntryParser();
        Iterable<?> programListItems = programEntryParser.getElements(channel, day);

        // pre-process items
        try {
            for (Iterator<?> it = programListItems.iterator(); it.hasNext(); ) {
                ProgramEntry programEntry = programEntryParser.preProcessItem(it.next(), day);

                if (programEntry == null) {
                    LOG.warning("Null program entry found. Skipping storage step.");
                    continue;
                }

                if (programEntry.getTechnicalId() == null || programEntry.getTechnicalId().isEmpty()) {
                    throw new ParserException("no technical id set by pre-process adapter method");
                }
                if (programEntryRepository.existsByTechnicalId(programEntry.getTechnicalId())) {
                    // update existing entry. avoid duplicates!
                    programEntry = programEntryRepository.getByTechnicalIdAndChannel(programEntry.getTechnicalId(), channel);
                } else {
                    initNewProgramEntry(programEntry, channel);
                }
                linkedProgramList.add(programEntry);
                entryCounter++;
            }
        } catch (ProgramEntryParserException e) {
            errorCounter++;
            LOG.throwing(AbstractReader.class.getName(), "handleProgramEntries", e);
        }

        programEntryParser.preProcessProgramList(linkedProgramList);

        LOG.info(String.format("Found %d entries, start post-processing", entryCounter));

        LOG.info(String.format("Successfully parsed %d entries. Errors: %d",
                linkedProgramList.size(), errorCounter));
        return linkedProgramList;
    }

    /**
     * method to init a new program entry
     *
     * @param programEntry
     * @param channel
     */
    public void initNewProgramEntry(ProgramEntry programEntry, Channel channel) {
        programEntry.setCreatedAt(LocalDateTime.now());
        // link channel record
        programEntry.setChannel(channel);
        programEntry.setAdapterFamily(getAdapterFamily());
    }

    /**
     * returns a filtered list of tv shows. only returning records which need to be updated.
     *
     * @return
     */
    private List<TvShow> detectTvShows() {
        int entryCounter = 0;
        int errorCounter = 0;
        List<TvShow> linkedTvList = new ArrayList<>();

        final TvShowParserInterface tvShowParser = getTvShowParser();
        try {
            List<TvShow> tvShowList = tvShowParser.getEntries();
            if (tvShowList == null) {
                throw new IllegalStateException("null tv show list given, use empty list.");
            }

            if (tvShowList.isEmpty()) {
                LOG.warning("empty tv show list given");
            }
            for (TvShow singleShow : tvShowList) {
                Optional<TvShow> dbRecordOptional = tvShowRepository
                        .findByTechnicalIdAndAdapterFamily(singleShow.getTechnicalId(), singleShow.getAdapterFamily());
                if (dbRecordOptional.isPresent()) {
                    // replace with entry from list
                    final TvShow tvShowDb = dbRecordOptional.get();
                    linkedTvList.add(tvShowDb);
                } else {
                    singleShow.setCreatedAt(LocalDateTime.now());
                    linkedTvList.add(singleShow);
                }
                entryCounter++;
            }
        } catch (TvShowParserException e) {
            errorCounter++;
            LOG.throwing(AbstractReader.class.getName(), "handleProgramEntries", e);
        }

        LOG.info(String.format("Found %d entries, start post-processing", entryCounter));

        LOG.info(String.format("Successfully parsed %d entries. Errors: %d",
                linkedTvList.size(), errorCounter));
        return linkedTvList;
    }

    /**
     * central update logic of stored entries
     *
     * @param model abstract record
     * @return whether to update this record or not - depends on preferences
     */
    protected boolean needsUpdate(UpdateControlInterface model) {
        if (!updateMode) {
            return false;
        }
        if (updateModeForce) {
            return true;
        }
        LocalDateTime lastUpdate;
        if (model.getUpdatedAt() == null) {
            lastUpdate = model.getCreatedAt();
        } else {
            lastUpdate = model.getUpdatedAt();
        }
        long minutes = ChronoUnit.MINUTES.between(lastUpdate, LocalDateTime.now());

        if (minutes > 60 * getInvalidationHours()) {
            return true;
        }
        return false;
    }

    /**
     * method to register a custom parser per adapter family
     *
     * @param customerParserClass
     */
    protected void registerParser(Class<? extends CustomParser> customerParserClass) {
        if (customerParserClass == null) {
            throw new IllegalArgumentException("null parser class given");
        }
        // fetch bean of registered parser class
        final CustomParser parserBean = context.getBean(customerParserClass);
        if (!this.customParsers.contains(parserBean)) {
            this.customParsers.add(parserBean);
        } else {
            LOG.warning(String.format("parser '%s' already registered", customerParserClass.getName()));
        }
    }

    public void setApplicationContext(ApplicationContext applicationContext) {
        this.context = applicationContext;
    }

    /**
     * method to retrieve channels of this adapter
     *
     * @return all channels the app knows
     */
    protected List<Channel> getAllChannels() {
        return channelService.getAllChannelsByFamily(getAdapterFamily());
    }

    private int getCollectDaysInFutureMax() {
        try {
            return Double.valueOf(collectDaysInFutureMax).intValue();
        } catch (NumberFormatException ignored) {
        }
        return 0;
    }

    private int getCollectDaysInPastMax() {
        try {
            return Double.valueOf(collectDaysInPastMax).intValue();
        } catch (NumberFormatException ignored) {
        }
        return 0;
    }

    private int getInvalidationHours() {
        try {
            return Double.valueOf(invalidationHours).intValue();
        } catch (NumberFormatException ignored) {
        }
        return 0;
    }

    private Optional<LocalDate> getStartDate() {
        return getLocalDate(startDate);
    }

    private Optional<LocalDate> getEndDate() {
        return getLocalDate(endDate);
    }

    /**
     * value to localdate optional
     *
     * @param dateValue iso date string
     * @return
     */
    private Optional<LocalDate> getLocalDate(String dateValue) {
        if (dateValue == null || dateValue.isEmpty()) {
            return Optional.empty();
        }
        try {
            return Optional.of(LocalDate.parse(dateValue));
        } catch (DateTimeParseException ignored) {
            LOG.warning("invalid end date value: " + dateValue);
        }
        return Optional.empty();
    }

    /**
     * this thread is used for tv shows
     */
    class TvShowPostProcessThread implements Runnable {
        private final TvShow tvShow;

        public TvShowPostProcessThread(TvShow tvShow) {
            this.tvShow = tvShow;
        }

        @Override
        @SuppressWarnings("squid:S1148")
        public void run() {
            try {
                LOG.finest("Start post-processing tv-show #" + this.tvShow.getId());
                getTvShowParser().postProcessEntry(this.tvShow);

                if (this.tvShow.getId() != null) {
                    this.tvShow.setUpdatedAt(LocalDateTime.now());
                }
                tvShowRepository.save(this.tvShow);
            } catch (TvShowParserException e) {
                LOG.throwing(TvShowPostProcessThread.class.getName(), "run", e);
            } catch (Exception e) {
                LOG.warning("tv show could not be stored: " + e.getMessage());
                LOG.throwing(TvShowPostProcessThread.class.getName(), "run", e);
                if (isDebug) {
                    LOG.info("Debug mode stacktrace: ");
                    e.printStackTrace();
                }
            } finally {
                LOG.finest("Finish post-processing tv-show #" + this.tvShow.getId());
            }
        }
    }
}

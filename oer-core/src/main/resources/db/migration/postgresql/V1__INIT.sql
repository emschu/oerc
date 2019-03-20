                 --
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.12
-- Dumped by pg_dump version 9.6.12

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: oer_server_dev; Type: DATABASE; Schema: -; Owner: postgres
--

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: channel; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.channel (
    id bigint NOT NULL,
    adapter_family integer,
    channel_key integer,
    home_page character varying(255),
    name character varying(255),
    technical_id character varying(255)
);


ALTER TABLE public.channel OWNER TO postgres;

--
-- Name: channel_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.channel_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.channel_id_seq OWNER TO postgres;

--
-- Name: channel_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.channel_id_seq OWNED BY public.channel.id;


--
-- Name: image_link; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.image_link (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    url character varying(2500)
);


ALTER TABLE public.image_link OWNER TO postgres;

--
-- Name: image_link_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.image_link_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.image_link_id_seq OWNER TO postgres;

--
-- Name: image_link_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.image_link_id_seq OWNED BY public.image_link.id;


--
-- Name: program_entry; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.program_entry (
    id bigint NOT NULL,
    adapter_family integer NOT NULL,
    created_at timestamp without time zone,
    description text,
    duration_in_minutes integer,
    end_date_time timestamp without time zone,
    home_page character varying(1000),
    start_date_time timestamp without time zone,
    technical_id character varying(500),
    title character varying(1000),
    updated_at timestamp without time zone,
    url character varying(1000),
    channel_id bigint NOT NULL
);


ALTER TABLE public.program_entry OWNER TO postgres;

--
-- Name: program_entry_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.program_entry_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.program_entry_id_seq OWNER TO postgres;

--
-- Name: program_entry_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.program_entry_id_seq OWNED BY public.program_entry.id;


--
-- Name: program_entry_image_links; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.program_entry_image_links (
    program_entry_id bigint NOT NULL,
    image_links_id bigint NOT NULL
);


ALTER TABLE public.program_entry_image_links OWNER TO postgres;

--
-- Name: program_entry_tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.program_entry_tags (
    program_entry_id bigint NOT NULL,
    tags_id bigint NOT NULL
);


ALTER TABLE public.program_entry_tags OWNER TO postgres;

--
-- Name: tag; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tag (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    tag_name character varying(500)
);


ALTER TABLE public.tag OWNER TO postgres;

--
-- Name: tag_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tag_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tag_id_seq OWNER TO postgres;

--
-- Name: tag_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tag_id_seq OWNED BY public.tag.id;


--
-- Name: tv_show; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tv_show (
    id bigint NOT NULL,
    adapter_family integer NOT NULL,
    additional_id character varying(1000),
    created_at timestamp without time zone,
    home_page character varying(1500),
    image_url character varying(1500),
    technical_id character varying(32) NOT NULL,
    title character varying(1000) NOT NULL,
    updated_at timestamp without time zone,
    url character varying(1500)
);


ALTER TABLE public.tv_show OWNER TO postgres;

--
-- Name: tv_show_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tv_show_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tv_show_id_seq OWNER TO postgres;

--
-- Name: tv_show_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tv_show_id_seq OWNED BY public.tv_show.id;


--
-- Name: tv_show_related_program_entries; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tv_show_related_program_entries (
    tv_show_id bigint NOT NULL,
    related_program_entry_id bigint NOT NULL
);


ALTER TABLE public.tv_show_related_program_entries OWNER TO postgres;

--
-- Name: tv_show_tags; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tv_show_tags (
    tv_show_id bigint NOT NULL,
    tag_id bigint NOT NULL
);


ALTER TABLE public.tv_show_tags OWNER TO postgres;

--
-- Name: channel id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.channel ALTER COLUMN id SET DEFAULT nextval('public.channel_id_seq'::regclass);


--
-- Name: image_link id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_link ALTER COLUMN id SET DEFAULT nextval('public.image_link_id_seq'::regclass);


--
-- Name: program_entry id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.program_entry ALTER COLUMN id SET DEFAULT nextval('public.program_entry_id_seq'::regclass);


--
-- Name: tag id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tag ALTER COLUMN id SET DEFAULT nextval('public.tag_id_seq'::regclass);


--
-- Name: tv_show id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tv_show ALTER COLUMN id SET DEFAULT nextval('public.tv_show_id_seq'::regclass);


--
-- Name: channel channel_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.channel
    ADD CONSTRAINT channel_pkey PRIMARY KEY (id);


--
-- Name: image_link image_link_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.image_link
    ADD CONSTRAINT image_link_pkey PRIMARY KEY (id);


--
-- Name: program_entry program_entry_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.program_entry
    ADD CONSTRAINT program_entry_pkey PRIMARY KEY (id);


--
-- Name: tag tag_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tag
    ADD CONSTRAINT tag_pkey PRIMARY KEY (id);


--
-- Name: tv_show tv_show_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tv_show
    ADD CONSTRAINT tv_show_pkey PRIMARY KEY (id);


--
-- Name: program_entry uk1yvibxe4i3tba0fvnrh62o23n; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.program_entry
    ADD CONSTRAINT uk1yvibxe4i3tba0fvnrh62o23n UNIQUE (technical_id, adapter_family);


--
-- Name: tag uk_1r1tyf6uga9k6jwdqnoqwtk2a; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tag
    ADD CONSTRAINT uk_1r1tyf6uga9k6jwdqnoqwtk2a UNIQUE (tag_name);


--
-- Name: program_entry uk_5q3jqqjravv11mdtot80qlb1m; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.program_entry
    ADD CONSTRAINT uk_5q3jqqjravv11mdtot80qlb1m UNIQUE (technical_id);


--
-- Name: channel uk_e7rgsr1tb3kbjo50ullunoo6p; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.channel
    ADD CONSTRAINT uk_e7rgsr1tb3kbjo50ullunoo6p UNIQUE (channel_key);


--
-- Name: channel uk_s25bgs344jgesyfn3mp89nhm4; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.channel
    ADD CONSTRAINT uk_s25bgs344jgesyfn3mp89nhm4 UNIQUE (technical_id);


--
-- Name: tv_show ukn7n8qfppeuky4u7tp0ahmau9h; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tv_show
    ADD CONSTRAINT ukn7n8qfppeuky4u7tp0ahmau9h UNIQUE (adapter_family, technical_id);


--
-- Name: program_entry_tags fk29tmop492x2w7lrux944oc8vr; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.program_entry_tags
    ADD CONSTRAINT fk29tmop492x2w7lrux944oc8vr FOREIGN KEY (tags_id) REFERENCES public.tag(id);


--
-- Name: program_entry fk5j0cawsql77ko5diak2uypvws; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.program_entry
    ADD CONSTRAINT fk5j0cawsql77ko5diak2uypvws FOREIGN KEY (channel_id) REFERENCES public.channel(id);


--
-- Name: program_entry_image_links fk6nno5x4eo2ivf29iljq3xv9v; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.program_entry_image_links
    ADD CONSTRAINT fk6nno5x4eo2ivf29iljq3xv9v FOREIGN KEY (program_entry_id) REFERENCES public.program_entry(id);


--
-- Name: tv_show_tags fk6xqwqofxo16ox0jd10y7ys9k8; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tv_show_tags
    ADD CONSTRAINT fk6xqwqofxo16ox0jd10y7ys9k8 FOREIGN KEY (tv_show_id) REFERENCES public.tv_show(id);


--
-- Name: tv_show_related_program_entries fkdjh011hl4e7x32o9oecj68mdd; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tv_show_related_program_entries
    ADD CONSTRAINT fkdjh011hl4e7x32o9oecj68mdd FOREIGN KEY (related_program_entry_id) REFERENCES public.program_entry(id);


--
-- Name: program_entry_tags fkhs976p4o50riioufcqb1eawx5; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.program_entry_tags
    ADD CONSTRAINT fkhs976p4o50riioufcqb1eawx5 FOREIGN KEY (program_entry_id) REFERENCES public.program_entry(id);


--
-- Name: tv_show_tags fkje6ehg0ldgyc1c07yimikwydv; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tv_show_tags
    ADD CONSTRAINT fkje6ehg0ldgyc1c07yimikwydv FOREIGN KEY (tag_id) REFERENCES public.tag(id);


--
-- Name: tv_show_related_program_entries fkk602gv4wskb2m2l45s72lns8g; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tv_show_related_program_entries
    ADD CONSTRAINT fkk602gv4wskb2m2l45s72lns8g FOREIGN KEY (tv_show_id) REFERENCES public.tv_show(id);


--
-- Name: program_entry_image_links fkshww5u9h6ei2qtxdo6yqobfp; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.program_entry_image_links
    ADD CONSTRAINT fkshww5u9h6ei2qtxdo6yqobfp FOREIGN KEY (image_links_id) REFERENCES public.image_link(id);


--
-- PostgreSQL database dump complete
--
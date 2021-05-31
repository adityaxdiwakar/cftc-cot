CREATE TABLE public.financial (
    id character varying NOT NULL,
    product character varying,
    date integer,
    deallong integer,
    dealshort integer,
    assetlong integer,
    assetshort integer,
    levlong integer,
    levshort integer,
    otherlong integer,
    othershort integer
);


ALTER TABLE ONLY public.financial
    ADD CONSTRAINT financial_pkey PRIMARY KEY (id);


CREATE TABLE public.financialprods (
    id character varying NOT NULL,
    name character varying
);

ALTER TABLE ONLY public.financialprods
    ADD CONSTRAINT financialprods_pkey PRIMARY KEY (id);


CREATE TABLE public.disaggregated (
    id character varying NOT NULL,
    product character varying,
    date integer,
    prodlong integer,
    prodshort integer,
    swaplong integer,
    swapshort integer,
    mmlong integer,
    mmshort integer,
    otherlong integer,
    othershort integer
);

ALTER TABLE ONLY public.disaggregated
    ADD CONSTRAINT disaggregated_pkey PRIMARY KEY (id);

CREATE TABLE public.disaggregatedprods (
    id character varying NOT NULL,
    name character varying
);


ALTER TABLE ONLY public.disaggregatedprods
    ADD CONSTRAINT disaggregatedprods_pkey PRIMARY KEY (id);


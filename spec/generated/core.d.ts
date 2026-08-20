export type BusinessEntityId = string;
export type EstablishmentId = string;
export type ProcessId = string;
export type ActivityId = string;
export type ProductOfferingId = string;
export type ChemicalHoldingId = string;
export type EquipmentItemId = string;
export type AttributeDefinitionUri = string;
export type FactId = string;
export type SchemeId = string;
export type SchemeVersionId = string;
export type ConceptId = string;
export type CorrespondenceId = string;
export type ConceptMappingId = string;
export type ClassificationAssignmentId = string;
export type CodeTranslationId = string;
export type RollupRuleId = string;
export type JurisdictionId = string;
export type JurisdictionProfileId = string;
export type AuthorityId = string;
export type LegalSourceId = string;
export type RegimeId = string;
export type ObligationId = string;
export type DeterminationId = string;
export type ProvenanceRecordId = string;
export type CredentialId = string;
export type RequirementId = string;
export type CredentialDependencyEdgeId = string;
/**
* The unit of analysis a fact, predicate, or obligation attaches to. Declared explicitly on every attribute and every predicate, because regulations key to different units and conflating them produces wrong determinations in both directions.
*/
export enum ScopeUnit {
    
    /** The legal entity as a whole. Company-wide headcount lives here. */
    business = "business",
    /** A single physical site. Most industry-code keying happens here. */
    establishment = "establishment",
    /** An activity unit within an establishment. Required by chemical-process regimes, which attach codes and thresholds below the site level. */
    process = "process",
    /** A line of business carried on at an establishment. */
    activity = "activity",
};
/**
* Disclosure sensitivity, which governs storage and model egress.
*/
export enum DataClass {
    
    public = "public",
    business_confidential = "business_confidential",
    /** Chemical inventories and comparable hazard data. Never transmitted to a third-party model provider or analytics endpoint. */
    restricted = "restricted",
};
/**
* Whether a classification has been affirmed by a person. Only a non-unconfirmed assignment may be referenced by a determination.
*/
export enum ConfirmationState {
    
    unconfirmed = "unconfirmed",
    owner_confirmed = "owner_confirmed",
    staff_confirmed = "staff_confirmed",
    agency_assigned = "agency_assigned",
};
/**
* How a candidate code came to be proposed or set.
*/
export enum AssignmentMethod {
    
    llm_candidate = "llm_candidate",
    embedding_retrieval = "embedding_retrieval",
    owner_selected = "owner_selected",
    agency_assigned = "agency_assigned",
    derived_rollup = "derived_rollup",
    crosswalk_translated = "crosswalk_translated",
};
/**
* A concept's fate across scheme revisions. `reused` is the dangerous one: the same code string carrying a different concept in a later vintage.
*/
export enum RevisionStatus {
    
    new = "new",
    revised = "revised",
    reused = "reused",
    retired = "retired",
    unchanged = "unchanged",
};
/**
* SKOS mapping relation between two concepts.
*/
export enum MatchType {
    
    exactMatch = "exactMatch",
    closeMatch = "closeMatch",
    broadMatch = "broadMatch",
    narrowMatch = "narrowMatch",
    relatedMatch = "relatedMatch",
};
/**
* Whether a code list is exhaustive. The single most consequential distinction in the model: a hit on either kind proves inclusion, but only a closed list can prove exclusion.
*/
export enum ListSemantics {
    
    /** The list is exhaustive. A predicate over it may return TRUE or FALSE. */
    enumerative_closed = "enumerative_closed",
    /** The list gives examples. A predicate over it may return TRUE or UNKNOWN, and never FALSE. Negation over such a list is a lint error. */
    illustrative_open = "illustrative_open",
};
/**
* Kleene three-valued logic, with the unknown case split by what resolves it.
*/
export enum TruthValue {
    
    TRUE = "TRUE",
    FALSE = "FALSE",
    /** Resolvable by asking the applicant a question. */
    UNKNOWN_MISSING_INPUT = "UNKNOWN_MISSING_INPUT",
    /** Not resolvable by the applicant at all; requires a determination by the issuing authority. */
    UNKNOWN_NOT_SELF_DETERMINABLE = "UNKNOWN_NOT_SELF_DETERMINABLE",
};
/**
* Boolean structure of a requirement or predicate tree node.
*/
export enum NodeType {
    
    AND_GROUP = "AND_GROUP",
    OR_GROUP = "OR_GROUP",
    LEAF = "LEAF",
};
/**
* How a code at one scope is derived from codes at a narrower scope.
*/
export enum RollupMethod {
    
    largest_receipts = "largest_receipts",
    largest_payroll = "largest_payroll",
    operator_designated = "operator_designated",
    union_all = "union_all",
    /** Roll-up is prohibited for this regime. Some programmes explicitly forbid defaulting a process-level code from the facility's primary code. */
    none = "none",
};

export enum JurisdictionLevel {
    
    nation = "nation",
    state = "state",
    county = "county",
    place = "place",
};

export enum CredentialType {
    
    license = "license",
    certification = "certification",
    registration = "registration",
    /** Local extension. CTDL has no permit class. */
    permit = "permit",
};


/**
 * A legal entity. The company scope.
 */
export interface BusinessEntity {
    id: string,
    legal_name: string,
    /** LLC, corporation, sole proprietorship, partnership. */
    entity_form?: string,
    formation_jurisdiction?: JurisdictionId,
    formation_date?: date,
    /** Peak employment across the ENTIRE company in the prior calendar year. Deliberately named with a company scope, because size-based exemptions are company-wide while industry-based exemptions are per-site. */
    company_peak_employment_prior_cy?: number,
    company_annual_receipts?: string,
    /** Affiliate grouping, which size standards aggregate across. */
    affiliation_group_id?: string,
    establishments?: Establishment[],
}


/**
 * A single physical site operated by a business entity.
 */
export interface Establishment {
    id: string,
    business_entity: BusinessEntityId,
    site_address?: string,
    /** Ordered nation to place. A set, never a scalar. */
    jurisdiction_path?: JurisdictionId[],
    is_unincorporated?: boolean,
    employment_at_establishment?: number,
    floor_area_sqft?: number,
    opened_on?: date,
    closed_on?: date,
    activities?: Activity[],
    processes?: Process[],
    product_offerings?: ProductOffering[],
}


/**
 * An activity unit beneath an establishment. Exists because chemical-process regimes attach codes and quantity thresholds below the site level, and forbid inheriting the site's primary code as a default.
 */
export interface Process {
    id: string,
    establishment: EstablishmentId,
    name?: string,
    process_type?: string,
    is_covered_process?: boolean,
    chemical_holdings?: ChemicalHolding[],
    equipment?: EquipmentItem[],
}


/**
 * A line of business carried on at an establishment.
 */
export interface Activity {
    id: string,
    establishment: EstablishmentId,
    description: string,
    /** Deliberately NOT uniquely constrained. Some programmes speak of "primary industrial activity(ies)" in the plural, and an activity matching a narrative category is primary regardless of its code. */
    is_primary?: boolean,
    /** receipts, headcount, production_rate, operator_designated, narrative_category */
    primacy_basis?: string,
    receipts_share?: string,
    /** Receipts ordering is stored as advisory guidance with an operator-supplied justification, never as an enforced rule. */
    primacy_justification_text?: string,
}


/**
 * A product or service offered. A separate demand-based axis; not derivable from the supply-based industry classification.
 */
export interface ProductOffering {
    id: string,
    establishment?: EstablishmentId,
    revenue_share?: string,
    is_regulated_substance?: boolean,
}


/**
 * A chemical held at a process. Data class is restricted.
 */
export interface ChemicalHolding {
    id: string,
    process: ProcessId,
    cas_number?: string,
    max_quantity?: string,
    uom?: string,
    flashpoint_c?: string,
    hazard_category?: string,
    container_type?: string,
    /** owner_declared, sds_parsed, or inspection. */
    source?: string,
}



export interface EquipmentItem {
    id: string,
    process?: ProcessId,
    equipment_type?: string,
    design_capacity?: string,
    uom?: string,
}


/**
 * The registry entry for a fact the system may collect or a rule may test. A rule referencing an unregistered attribute, or one at a scope that disagrees with this entry, fails CI. This is what makes scope-confusion bugs unwriteable.
 */
export interface AttributeDefinition {
    uri: string,
    label: string,
    scope_unit: string,
    datatype: string,
    unit_dimension?: string,
    enum_ref?: string,
    collection_method?: string,
    data_class: string,
    llm_egress_allowed: boolean,
}


/**
 * An open-world fact carrying the long tail of attributes that do not warrant a typed column.
 */
export interface Fact {
    id: string,
    subject_ref: string,
    subject_scope: string,
    attribute: AttributeDefinitionUri,
    value_typed: string,
    unit?: string,
    effective_from?: date,
    effective_to?: date,
    source?: string,
    confidence?: string,
}


/**
 * A classification scheme.
 */
export interface Scheme {
    id: string,
    label: string,
    publisher?: string,
    /** False for code spaces that look like industry codes but are not, and must never be crosswalked as though they were. */
    is_industry_scheme?: boolean,
}


/**
 * A dated edition of a scheme. Vintage is part of concept identity.
 */
export interface SchemeVersion {
    id: string,
    scheme: SchemeId,
    vintage: string,
    effective_from?: date,
    effective_to?: date,
}


/**
 * A single code within a scheme version. Identity is the triple (scheme, vintage, code) -- never a bare code.
 */
export interface Concept {
    /** Synthetic key. The natural key is (scheme, vintage, code). */
    id: string,
    scheme_version: SchemeVersionId,
    code: string,
    title: string,
    definition_text?: string,
    level?: number,
    parent_code?: string,
    revision_status?: string,
}


/**
 * A published crosswalk between two scheme versions, as an addressable and versionable object rather than a lookup table.
 */
export interface Correspondence {
    id: string,
    source_scheme_version: SchemeVersionId,
    target_scheme_version: SchemeVersionId,
    publisher?: string,
    published_on?: date,
    file_provenance?: string,
    coverage_notes?: string,
}


/**
 * One row of a correspondence.
 */
export interface ConceptMapping {
    id: string,
    correspondence: CorrespondenceId,
    source_concept?: ConceptId,
    target_concept?: ConceptId,
    match_type: string,
    /** Local extension. XKOS 1.2 defines no such property. */
    apportionment_ratio?: string,
    /** Local extension. */
    match_strength?: string,
    asserted_by?: string,
}


/**
 * A candidate or confirmed code for a subject. Candidates are presented ranked for human confirmation; only a confirmed assignment may be the code of record.
 */
export interface ClassificationAssignment {
    id: string,
    subject_ref: string,
    subject_scope: string,
    scheme_version: SchemeVersionId,
    code: string,
    rank?: number,
    score?: string,
    method: string,
    confirmation_state: string,
    is_code_of_record: boolean,
    /** What the applicant was offered, not merely what they chose. Competent people disagree on the same description; an appeal needs the offered set. */
    alternatives_shown?: string[],
    provenance?: ProvenanceRecordId,
}


/**
 * The result of walking one or more correspondences. The hop path is retained so a translation can be audited rather than trusted.
 */
export interface CodeTranslation {
    id: string,
    assignment: ClassificationAssignmentId,
    target_scheme_version: SchemeVersionId,
    result_code?: string[],
    /** Ordered correspondence identifiers. */
    hop_path?: string[],
    match_type_chain?: string,
    /** True only when every hop is an exactMatch. Close-match chains never auto-compose. */
    is_composable?: boolean,
    review_required?: boolean,
}


/**
 * How a code at a wider scope derives from a narrower one, as data. Encoding this as rows is what lets one regime forbid the roll-up that another requires.
 */
export interface RollupRule {
    id: string,
    target_scope: string,
    source_scope: string,
    method: string,
    authority_citation: string,
    regime?: RegimeId,
    no_default_from_parent?: boolean,
}



export interface Jurisdiction {
    id: string,
    level: string,
    geoid?: string,
    label: string,
    parent?: JurisdictionId,
    effective_from?: date,
    effective_to?: date,
}


/**
 * A jurisdiction expressed as an inclusion set with exceptions, because "everywhere except one state" is a common real pattern that a flat column cannot represent.
 */
export interface JurisdictionProfile {
    id: string,
    main_jurisdiction?: JurisdictionId[],
    jurisdiction_exception?: JurisdictionId[],
    global_flag?: boolean,
    residency_required?: boolean,
}


/**
 * An agency that issues credentials or administers obligations.
 */
export interface Authority {
    id: string,
    label: string,
    jurisdiction?: JurisdictionId,
}


/**
 * A retrieved, hashed, point-in-time snapshot of authoritative text. Every obligation traces to one of these.
 */
export interface LegalSource {
    id: string,
    citation: string,
    source_url: string,
    source_system?: string,
    text_hash?: string,
    retrieved_at?: string,
    as_of_date?: date,
    amendment_date?: date,
}


/**
 * A regulatory programme that generates obligations.
 */
export interface Regime {
    id: string,
    label: string,
    scope_unit: string,
    authority?: AuthorityId,
}



export interface Obligation {
    id: string,
    regime: RegimeId,
    obligation_type: string,
    scope_unit: string,
    trigger_rule_id: string,
    legal_source: LegalSourceId,
    deadline_rule?: string,
    recurrence?: string,
    /** Survives every exemption in its regime. A non-waivable obligation is surfaced even when the applicant is otherwise exempt. */
    non_waivable?: boolean,
}


/**
 * A reproducible answer. Pinning the rule version, the law date, and a hash of the inputs is what lets the system answer "why did it say that in March".
 */
export interface Determination {
    id: string,
    subject_ref: string,
    subject_scope: string,
    obligation: ObligationId,
    result: string,
    /** MUST reference an assignment whose confirmation_state is not `unconfirmed`. Enforced by a database CHECK constraint, not by application code. */
    classification_assignment?: ClassificationAssignmentId,
    rule_version_id: string,
    engine_version: string,
    as_of_law: date,
    input_snapshot_hash: string,
    /** What must be answered next when the result is indeterminate. */
    missing_attributes?: AttributeDefinitionUri[],
    /** The predicate tree annotated with truth values and citations. */
    evidence_tree?: string,
    determined_at?: string,
}


/**
 * What produced a surfaced item. No obligation is displayed without one.
 */
export interface ProvenanceRecord {
    id: string,
    model_id?: string,
    prompt_version?: string,
    index_version?: string,
    rank_at_selection?: number,
    confirmed_by?: string,
    confirmed_at?: string,
    ui_version?: string,
}


/**
 * A licence, certification, registration, or permit.
 */
export interface Credential {
    id: string,
    credential_type: string,
    label: string,
    issuing_authority: AuthorityId,
    jurisdiction_profile?: JurisdictionProfileId,
    /** Adopted from CTDL, which types this as an untyped string with no vintage validation, so codes from different revisions are indistinguishable. */
    industry_codes?: string[],
    /** Local extension supplying the vintage CTDL omits. */
    industry_code_vintage?: string,
    legal_source?: LegalSourceId,
    estimated_cost?: string,
    renewal_period_months?: number,
}


/**
 * A reified condition between a credential and what it demands. A table rather than a foreign key, because the edge carries jurisdiction, residency, experience, cost, and dates. Recursive: AND groups contain OR alternative sets contain further groups, which is what real licensing paths need.
 */
export interface Requirement {
    id: string,
    credential: CredentialId,
    parent?: RequirementId,
    node_type: string,
    /** prerequisite, concurrent, conditional-on, or renewal. */
    edge_kind?: string,
    jurisdiction_profile?: JurisdictionProfileId,
    residency_required?: boolean,
    min_age?: number,
    years_experience?: string,
    estimated_cost?: string,
    effective_from?: date,
    effective_to?: date,
    legal_source: LegalSourceId,
    target_credential?: CredentialId,
    /** A predicate from the rules DSL. This is the seam that unifies the regulatory and credential halves of the system into a single evaluator. */
    target_predicate?: string,
}


/**
 * A derived, typed edge used for topological ordering. Cycles are a transcription error in the source law and fail CI; they are not a runtime condition to handle.
 */
export interface CredentialDependencyEdge {
    id: string,
    from_credential: CredentialId,
    to_credential: CredentialId,
    edge_kind: string,
    derived_from_requirement?: RequirementId,
}




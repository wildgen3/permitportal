from __future__ import annotations

import re
import sys
from datetime import (
    date,
    datetime,
    time
)
from decimal import Decimal
from enum import Enum
from typing import (
    Any,
    ClassVar,
    Literal,
    Optional,
    Union
)

from pydantic import (
    BaseModel,
    ConfigDict,
    Field,
    RootModel,
    SerializationInfo,
    SerializerFunctionWrapHandler,
    field_validator,
    model_serializer
)


metamodel_version = "1.11.0"
version = "None"


class ConfiguredBaseModel(BaseModel):
    model_config = ConfigDict(
        serialize_by_alias = True,
        validate_by_name = True,
        validate_assignment = True,
        validate_default = True,
        extra = "forbid",
        arbitrary_types_allowed = True,
        use_enum_values = True,
        strict = False,
    )





class LinkMLMeta(RootModel):
    root: dict[str, Any] = {}
    model_config = ConfigDict(frozen=True)

    def __getattr__(self, key:str):
        return getattr(self.root, key)

    def __getitem__(self, key:str):
        return self.root[key]

    def __setitem__(self, key:str, value):
        self.root[key] = value

    def __contains__(self, key:str) -> bool:
        return key in self.root


linkml_meta = LinkMLMeta({'default_prefix': 'pg',
     'default_range': 'string',
     'description': 'The canonical model of a business, its classification, the '
                    'obligations that attach to it, and the credentials those '
                    'obligations require.\n'
                    'Two invariants govern everything here. First, an industry '
                    'code is never a bare string: the identity of a concept is the '
                    'triple (scheme, vintage, code), because codes have been '
                    'reused for different concepts across revisions. Second, an '
                    'unconfirmed classification may never reach a determination -- '
                    'enforced at the database layer, not by application code.',
     'id': 'https://wildgen3.github.io/permitgraph/spec/model/core',
     'imports': ['linkml:types'],
     'license': 'https://creativecommons.org/licenses/by/4.0/',
     'name': 'permitgraph-core',
     'prefixes': {'ceterms': {'prefix_prefix': 'ceterms',
                              'prefix_reference': 'https://purl.org/ctdl/terms/'},
                  'dcterms': {'prefix_prefix': 'dcterms',
                              'prefix_reference': 'http://purl.org/dc/terms/'},
                  'linkml': {'prefix_prefix': 'linkml',
                             'prefix_reference': 'https://w3id.org/linkml/'},
                  'pg': {'prefix_prefix': 'pg',
                         'prefix_reference': 'https://wildgen3.github.io/permitgraph/spec/model/core/'},
                  'skos': {'prefix_prefix': 'skos',
                           'prefix_reference': 'http://www.w3.org/2004/02/skos/core#'},
                  'xkos': {'prefix_prefix': 'xkos',
                           'prefix_reference': 'http://rdf-vocabulary.ddialliance.org/xkos#'}},
     'source_file': 'spec/model/core.yaml',
     'title': 'PermitGraph canonical model'} )

class ScopeUnit(str, Enum):
    """
    The unit of analysis a fact, predicate, or obligation attaches to. Declared explicitly on every attribute and every predicate, because regulations key to different units and conflating them produces wrong determinations in both directions.
    """
    business = "business"
    """
    The legal entity as a whole. Company-wide headcount lives here.
    """
    establishment = "establishment"
    """
    A single physical site. Most industry-code keying happens here.
    """
    process = "process"
    """
    An activity unit within an establishment. Required by chemical-process regimes, which attach codes and thresholds below the site level.
    """
    activity = "activity"
    """
    A line of business carried on at an establishment.
    """


class DataClass(str, Enum):
    """
    Disclosure sensitivity, which governs storage and model egress.
    """
    public = "public"
    business_confidential = "business_confidential"
    restricted = "restricted"
    """
    Chemical inventories and comparable hazard data. Never transmitted to a third-party model provider or analytics endpoint.
    """


class ConfirmationState(str, Enum):
    """
    Whether a classification has been affirmed by a person. Only a non-unconfirmed assignment may be referenced by a determination.
    """
    unconfirmed = "unconfirmed"
    owner_confirmed = "owner_confirmed"
    staff_confirmed = "staff_confirmed"
    agency_assigned = "agency_assigned"


class AssignmentMethod(str, Enum):
    """
    How a candidate code came to be proposed or set.
    """
    llm_candidate = "llm_candidate"
    embedding_retrieval = "embedding_retrieval"
    owner_selected = "owner_selected"
    agency_assigned = "agency_assigned"
    derived_rollup = "derived_rollup"
    crosswalk_translated = "crosswalk_translated"


class RevisionStatus(str, Enum):
    """
    A concept's fate across scheme revisions. `reused` is the dangerous one: the same code string carrying a different concept in a later vintage.
    """
    new = "new"
    revised = "revised"
    reused = "reused"
    retired = "retired"
    unchanged = "unchanged"


class MatchType(str, Enum):
    """
    SKOS mapping relation between two concepts.
    """
    exactMatch = "exactMatch"
    closeMatch = "closeMatch"
    broadMatch = "broadMatch"
    narrowMatch = "narrowMatch"
    relatedMatch = "relatedMatch"


class ListSemantics(str, Enum):
    """
    Whether a code list is exhaustive. The single most consequential distinction in the model: a hit on either kind proves inclusion, but only a closed list can prove exclusion.
    """
    enumerative_closed = "enumerative_closed"
    """
    The list is exhaustive. A predicate over it may return TRUE or FALSE.
    """
    illustrative_open = "illustrative_open"
    """
    The list gives examples. A predicate over it may return TRUE or UNKNOWN, and never FALSE. Negation over such a list is a lint error.
    """


class TruthValue(str, Enum):
    """
    Kleene three-valued logic, with the unknown case split by what resolves it.
    """
    TRUE = "TRUE"
    FALSE = "FALSE"
    UNKNOWN_MISSING_INPUT = "UNKNOWN_MISSING_INPUT"
    """
    Resolvable by asking the applicant a question.
    """
    UNKNOWN_NOT_SELF_DETERMINABLE = "UNKNOWN_NOT_SELF_DETERMINABLE"
    """
    Not resolvable by the applicant at all; requires a determination by the issuing authority.
    """


class NodeType(str, Enum):
    """
    Boolean structure of a requirement or predicate tree node.
    """
    AND_GROUP = "AND_GROUP"
    OR_GROUP = "OR_GROUP"
    LEAF = "LEAF"


class RollupMethod(str, Enum):
    """
    How a code at one scope is derived from codes at a narrower scope.
    """
    largest_receipts = "largest_receipts"
    largest_payroll = "largest_payroll"
    operator_designated = "operator_designated"
    union_all = "union_all"
    none = "none"
    """
    Roll-up is prohibited for this regime. Some programmes explicitly forbid defaulting a process-level code from the facility's primary code.
    """


class JurisdictionLevel(str, Enum):
    nation = "nation"
    state = "state"
    county = "county"
    place = "place"


class CredentialType(str, Enum):
    license = "license"
    certification = "certification"
    registration = "registration"
    permit = "permit"
    """
    Local extension. CTDL has no permit class.
    """



class BusinessEntity(ConfiguredBaseModel):
    """
    A legal entity. The company scope.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    legal_name: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity']} })
    entity_form: Optional[str] = Field(default=None, description="""LLC, corporation, sole proprietorship, partnership.""", json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity']} })
    formation_jurisdiction: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity']} })
    formation_date: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity']} })
    company_peak_employment_prior_cy: Optional[int] = Field(default=None, description="""Peak employment across the ENTIRE company in the prior calendar year. Deliberately named with a company scope, because size-based exemptions are company-wide while industry-based exemptions are per-site.""", json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity']} })
    company_annual_receipts: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity']} })
    affiliation_group_id: Optional[str] = Field(default=None, description="""Affiliate grouping, which size standards aggregate across.""", json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity']} })
    establishments: Optional[list[Establishment]] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity']} })


class Establishment(ConfiguredBaseModel):
    """
    A single physical site operated by a business entity.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    business_entity: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })
    site_address: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })
    jurisdiction_path: Optional[list[str]] = Field(default=None, description="""Ordered nation to place. A set, never a scalar.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })
    is_unincorporated: Optional[bool] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })
    employment_at_establishment: Optional[int] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })
    floor_area_sqft: Optional[int] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })
    opened_on: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })
    closed_on: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })
    activities: Optional[list[Activity]] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })
    processes: Optional[list[Process]] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })
    product_offerings: Optional[list[ProductOffering]] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Establishment']} })


class Process(ConfiguredBaseModel):
    """
    An activity unit beneath an establishment. Exists because chemical-process regimes attach codes and quantity thresholds below the site level, and forbid inheriting the site's primary code as a default.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    establishment: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Process', 'Activity', 'ProductOffering']} })
    name: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Process']} })
    process_type: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Process']} })
    is_covered_process: Optional[bool] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Process']} })
    chemical_holdings: Optional[list[ChemicalHolding]] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Process']} })
    equipment: Optional[list[EquipmentItem]] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Process']} })


class Activity(ConfiguredBaseModel):
    """
    A line of business carried on at an establishment.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    establishment: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Process', 'Activity', 'ProductOffering']} })
    description: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Activity']} })
    is_primary: Optional[bool] = Field(default=None, description="""Deliberately NOT uniquely constrained. Some programmes speak of \"primary industrial activity(ies)\" in the plural, and an activity matching a narrative category is primary regardless of its code.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Activity']} })
    primacy_basis: Optional[str] = Field(default=None, description="""receipts, headcount, production_rate, operator_designated, narrative_category""", json_schema_extra = { "linkml_meta": {'domain_of': ['Activity']} })
    receipts_share: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Activity']} })
    primacy_justification_text: Optional[str] = Field(default=None, description="""Receipts ordering is stored as advisory guidance with an operator-supplied justification, never as an enforced rule.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Activity']} })


class ProductOffering(ConfiguredBaseModel):
    """
    A product or service offered. A separate demand-based axis; not derivable from the supply-based industry classification.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    establishment: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Process', 'Activity', 'ProductOffering']} })
    revenue_share: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ProductOffering']} })
    is_regulated_substance: Optional[bool] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ProductOffering']} })


class ChemicalHolding(ConfiguredBaseModel):
    """
    A chemical held at a process. Data class is restricted.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    process: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding', 'EquipmentItem']} })
    cas_number: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding']} })
    max_quantity: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding']} })
    uom: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding', 'EquipmentItem']} })
    flashpoint_c: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding']} })
    hazard_category: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding']} })
    container_type: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding']} })
    source: Optional[str] = Field(default=None, description="""owner_declared, sds_parsed, or inspection.""", json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding', 'Fact']} })


class EquipmentItem(ConfiguredBaseModel):
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    process: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding', 'EquipmentItem']} })
    equipment_type: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['EquipmentItem']} })
    design_capacity: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['EquipmentItem']} })
    uom: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding', 'EquipmentItem']} })


class AttributeDefinition(ConfiguredBaseModel):
    """
    The registry entry for a fact the system may collect or a rule may test. A rule referencing an unregistered attribute, or one at a scope that disagrees with this entry, fails CI. This is what makes scope-confusion bugs unwriteable.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    uri: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition']} })
    label: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition',
                       'Scheme',
                       'Jurisdiction',
                       'Authority',
                       'Regime',
                       'Credential']} })
    scope_unit: ScopeUnit = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition', 'Regime', 'Obligation']} })
    datatype: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition']} })
    unit_dimension: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition']} })
    enum_ref: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition']} })
    collection_method: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition']} })
    data_class: DataClass = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition']} })
    llm_egress_allowed: bool = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition']} })


class Fact(ConfiguredBaseModel):
    """
    An open-world fact carrying the long tail of attributes that do not warrant a typed column.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    subject_ref: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'ClassificationAssignment', 'Determination']} })
    subject_scope: ScopeUnit = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'ClassificationAssignment', 'Determination']} })
    attribute: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Fact']} })
    value_typed: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Fact']} })
    unit: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Fact']} })
    effective_from: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'SchemeVersion', 'Jurisdiction', 'Requirement']} })
    effective_to: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'SchemeVersion', 'Jurisdiction', 'Requirement']} })
    source: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ChemicalHolding', 'Fact']} })
    confidence: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Fact']} })


class Scheme(ConfiguredBaseModel):
    """
    A classification scheme.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    label: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition',
                       'Scheme',
                       'Jurisdiction',
                       'Authority',
                       'Regime',
                       'Credential']} })
    publisher: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Scheme', 'Correspondence']} })
    is_industry_scheme: Optional[bool] = Field(default=None, description="""False for code spaces that look like industry codes but are not, and must never be crosswalked as though they were.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Scheme']} })


class SchemeVersion(ConfiguredBaseModel):
    """
    A dated edition of a scheme. Vintage is part of concept identity.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    scheme: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['SchemeVersion']} })
    vintage: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['SchemeVersion']} })
    effective_from: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'SchemeVersion', 'Jurisdiction', 'Requirement']} })
    effective_to: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'SchemeVersion', 'Jurisdiction', 'Requirement']} })


class Concept(ConfiguredBaseModel):
    """
    A single code within a scheme version. Identity is the triple (scheme, vintage, code) -- never a bare code.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core',
         'unique_keys': {'natural_key': {'description': 'The real identity of a '
                                                        'concept.',
                                         'unique_key_name': 'natural_key',
                                         'unique_key_slots': ['scheme_version',
                                                              'code']}}})

    id: str = Field(default=..., description="""Synthetic key. The natural key is (scheme, vintage, code).""", json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    scheme_version: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Concept', 'ClassificationAssignment']} })
    code: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Concept', 'ClassificationAssignment']} })
    title: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Concept']} })
    definition_text: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Concept']} })
    level: Optional[int] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Concept', 'Jurisdiction']} })
    parent_code: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Concept']} })
    revision_status: Optional[RevisionStatus] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Concept']} })


class Correspondence(ConfiguredBaseModel):
    """
    A published crosswalk between two scheme versions, as an addressable and versionable object rather than a lookup table.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'exact_mappings': ['xkos:Correspondence'],
         'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    source_scheme_version: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Correspondence']} })
    target_scheme_version: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Correspondence', 'CodeTranslation']} })
    publisher: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Scheme', 'Correspondence']} })
    published_on: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Correspondence']} })
    file_provenance: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Correspondence']} })
    coverage_notes: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Correspondence']} })


class ConceptMapping(ConfiguredBaseModel):
    """
    One row of a correspondence.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    correspondence: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['ConceptMapping']} })
    source_concept: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ConceptMapping']} })
    target_concept: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ConceptMapping']} })
    match_type: MatchType = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['ConceptMapping']} })
    apportionment_ratio: Optional[Decimal] = Field(default=None, description="""Local extension. XKOS 1.2 defines no such property.""", json_schema_extra = { "linkml_meta": {'domain_of': ['ConceptMapping']} })
    match_strength: Optional[Decimal] = Field(default=None, description="""Local extension.""", json_schema_extra = { "linkml_meta": {'domain_of': ['ConceptMapping']} })
    asserted_by: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ConceptMapping']} })


class ClassificationAssignment(ConfiguredBaseModel):
    """
    A candidate or confirmed code for a subject. Candidates are presented ranked for human confirmation; only a confirmed assignment may be the code of record.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    subject_ref: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'ClassificationAssignment', 'Determination']} })
    subject_scope: ScopeUnit = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'ClassificationAssignment', 'Determination']} })
    scheme_version: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Concept', 'ClassificationAssignment']} })
    code: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Concept', 'ClassificationAssignment']} })
    rank: Optional[int] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ClassificationAssignment']} })
    score: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ClassificationAssignment']} })
    method: AssignmentMethod = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['ClassificationAssignment', 'RollupRule']} })
    confirmation_state: ConfirmationState = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['ClassificationAssignment']} })
    is_code_of_record: bool = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['ClassificationAssignment']} })
    alternatives_shown: Optional[list[str]] = Field(default=None, description="""What the applicant was offered, not merely what they chose. Competent people disagree on the same description; an appeal needs the offered set.""", json_schema_extra = { "linkml_meta": {'domain_of': ['ClassificationAssignment']} })
    provenance: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ClassificationAssignment']} })


class CodeTranslation(ConfiguredBaseModel):
    """
    The result of walking one or more correspondences. The hop path is retained so a translation can be audited rather than trusted.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    assignment: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['CodeTranslation']} })
    target_scheme_version: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Correspondence', 'CodeTranslation']} })
    result_code: Optional[list[str]] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['CodeTranslation']} })
    hop_path: Optional[list[str]] = Field(default=None, description="""Ordered correspondence identifiers.""", json_schema_extra = { "linkml_meta": {'domain_of': ['CodeTranslation']} })
    match_type_chain: Optional[list[MatchType]] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['CodeTranslation']} })
    is_composable: Optional[bool] = Field(default=None, description="""True only when every hop is an exactMatch. Close-match chains never auto-compose.""", json_schema_extra = { "linkml_meta": {'domain_of': ['CodeTranslation']} })
    review_required: Optional[bool] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['CodeTranslation']} })


class RollupRule(ConfiguredBaseModel):
    """
    How a code at a wider scope derives from a narrower one, as data. Encoding this as rows is what lets one regime forbid the roll-up that another requires.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    target_scope: ScopeUnit = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['RollupRule']} })
    source_scope: ScopeUnit = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['RollupRule']} })
    method: RollupMethod = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['ClassificationAssignment', 'RollupRule']} })
    authority_citation: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['RollupRule']} })
    regime: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['RollupRule', 'Obligation']} })
    no_default_from_parent: Optional[bool] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['RollupRule']} })


class Jurisdiction(ConfiguredBaseModel):
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    level: JurisdictionLevel = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Concept', 'Jurisdiction']} })
    geoid: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Jurisdiction']} })
    label: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition',
                       'Scheme',
                       'Jurisdiction',
                       'Authority',
                       'Regime',
                       'Credential']} })
    parent: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Jurisdiction', 'Requirement']} })
    effective_from: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'SchemeVersion', 'Jurisdiction', 'Requirement']} })
    effective_to: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'SchemeVersion', 'Jurisdiction', 'Requirement']} })


class JurisdictionProfile(ConfiguredBaseModel):
    """
    A jurisdiction expressed as an inclusion set with exceptions, because \"everywhere except one state\" is a common real pattern that a flat column cannot represent.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'exact_mappings': ['ceterms:JurisdictionProfile'],
         'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    main_jurisdiction: Optional[list[str]] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['JurisdictionProfile']} })
    jurisdiction_exception: Optional[list[str]] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['JurisdictionProfile']} })
    global_flag: Optional[bool] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['JurisdictionProfile']} })
    residency_required: Optional[bool] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['JurisdictionProfile', 'Requirement']} })


class Authority(ConfiguredBaseModel):
    """
    An agency that issues credentials or administers obligations.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    label: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition',
                       'Scheme',
                       'Jurisdiction',
                       'Authority',
                       'Regime',
                       'Credential']} })
    jurisdiction: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Authority']} })


class LegalSource(ConfiguredBaseModel):
    """
    A retrieved, hashed, point-in-time snapshot of authoritative text. Every obligation traces to one of these.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    citation: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['LegalSource']} })
    source_url: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['LegalSource']} })
    source_system: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['LegalSource']} })
    text_hash: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['LegalSource']} })
    retrieved_at: Optional[datetime ] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['LegalSource']} })
    as_of_date: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['LegalSource']} })
    amendment_date: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['LegalSource']} })


class Regime(ConfiguredBaseModel):
    """
    A regulatory programme that generates obligations.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    label: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition',
                       'Scheme',
                       'Jurisdiction',
                       'Authority',
                       'Regime',
                       'Credential']} })
    scope_unit: ScopeUnit = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition', 'Regime', 'Obligation']} })
    authority: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Regime']} })


class Obligation(ConfiguredBaseModel):
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    regime: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['RollupRule', 'Obligation']} })
    obligation_type: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Obligation']} })
    scope_unit: ScopeUnit = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition', 'Regime', 'Obligation']} })
    trigger_rule_id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Obligation']} })
    legal_source: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Obligation', 'Credential', 'Requirement']} })
    deadline_rule: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Obligation']} })
    recurrence: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Obligation']} })
    non_waivable: Optional[bool] = Field(default=None, description="""Survives every exemption in its regime. A non-waivable obligation is surfaced even when the applicant is otherwise exempt.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Obligation']} })


class Determination(ConfiguredBaseModel):
    """
    A reproducible answer. Pinning the rule version, the law date, and a hash of the inputs is what lets the system answer \"why did it say that in March\".
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    subject_ref: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'ClassificationAssignment', 'Determination']} })
    subject_scope: ScopeUnit = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'ClassificationAssignment', 'Determination']} })
    obligation: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Determination']} })
    result: TruthValue = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Determination']} })
    classification_assignment: Optional[str] = Field(default=None, description="""MUST reference an assignment whose confirmation_state is not `unconfirmed`. Enforced by a database CHECK constraint, not by application code.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Determination']} })
    rule_version_id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Determination']} })
    engine_version: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Determination']} })
    as_of_law: date = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Determination']} })
    input_snapshot_hash: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Determination']} })
    missing_attributes: Optional[list[str]] = Field(default=None, description="""What must be answered next when the result is indeterminate.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Determination']} })
    evidence_tree: Optional[str] = Field(default=None, description="""The predicate tree annotated with truth values and citations.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Determination']} })
    determined_at: Optional[datetime ] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Determination']} })


class ProvenanceRecord(ConfiguredBaseModel):
    """
    What produced a surfaced item. No obligation is displayed without one.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    model_id: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ProvenanceRecord']} })
    prompt_version: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ProvenanceRecord']} })
    index_version: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ProvenanceRecord']} })
    rank_at_selection: Optional[int] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ProvenanceRecord']} })
    confirmed_by: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ProvenanceRecord']} })
    confirmed_at: Optional[datetime ] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ProvenanceRecord']} })
    ui_version: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['ProvenanceRecord']} })


class Credential(ConfiguredBaseModel):
    """
    A licence, certification, registration, or permit.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    credential_type: CredentialType = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Credential']} })
    label: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['AttributeDefinition',
                       'Scheme',
                       'Jurisdiction',
                       'Authority',
                       'Regime',
                       'Credential']} })
    issuing_authority: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Credential']} })
    jurisdiction_profile: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Credential', 'Requirement']} })
    industry_codes: Optional[list[str]] = Field(default=None, description="""Adopted from CTDL, which types this as an untyped string with no vintage validation, so codes from different revisions are indistinguishable.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Credential']} })
    industry_code_vintage: Optional[str] = Field(default=None, description="""Local extension supplying the vintage CTDL omits.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Credential']} })
    legal_source: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Obligation', 'Credential', 'Requirement']} })
    estimated_cost: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Credential', 'Requirement']} })
    renewal_period_months: Optional[int] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Credential']} })


class Requirement(ConfiguredBaseModel):
    """
    A reified condition between a credential and what it demands. A table rather than a foreign key, because the edge carries jurisdiction, residency, experience, cost, and dates. Recursive: AND groups contain OR alternative sets contain further groups, which is what real licensing paths need.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'exact_mappings': ['ceterms:ConditionProfile'],
         'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    credential: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Requirement']} })
    parent: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Jurisdiction', 'Requirement']} })
    node_type: NodeType = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Requirement']} })
    edge_kind: Optional[str] = Field(default=None, description="""prerequisite, concurrent, conditional-on, or renewal.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Requirement', 'CredentialDependencyEdge']} })
    jurisdiction_profile: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Credential', 'Requirement']} })
    residency_required: Optional[bool] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['JurisdictionProfile', 'Requirement']} })
    min_age: Optional[int] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Requirement']} })
    years_experience: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Requirement']} })
    estimated_cost: Optional[Decimal] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Credential', 'Requirement']} })
    effective_from: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'SchemeVersion', 'Jurisdiction', 'Requirement']} })
    effective_to: Optional[date] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Fact', 'SchemeVersion', 'Jurisdiction', 'Requirement']} })
    legal_source: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Obligation', 'Credential', 'Requirement']} })
    target_credential: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['Requirement']} })
    target_predicate: Optional[str] = Field(default=None, description="""A predicate from the rules DSL. This is the seam that unifies the regulatory and credential halves of the system into a single evaluator.""", json_schema_extra = { "linkml_meta": {'domain_of': ['Requirement']} })


class CredentialDependencyEdge(ConfiguredBaseModel):
    """
    A derived, typed edge used for topological ordering. Cycles are a transcription error in the source law and fail CI; they are not a runtime condition to handle.
    """
    linkml_meta: ClassVar[LinkMLMeta] = LinkMLMeta({'from_schema': 'https://wildgen3.github.io/permitgraph/spec/model/core'})

    id: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['BusinessEntity',
                       'Establishment',
                       'Process',
                       'Activity',
                       'ProductOffering',
                       'ChemicalHolding',
                       'EquipmentItem',
                       'Fact',
                       'Scheme',
                       'SchemeVersion',
                       'Concept',
                       'Correspondence',
                       'ConceptMapping',
                       'ClassificationAssignment',
                       'CodeTranslation',
                       'RollupRule',
                       'Jurisdiction',
                       'JurisdictionProfile',
                       'Authority',
                       'LegalSource',
                       'Regime',
                       'Obligation',
                       'Determination',
                       'ProvenanceRecord',
                       'Credential',
                       'Requirement',
                       'CredentialDependencyEdge']} })
    from_credential: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['CredentialDependencyEdge']} })
    to_credential: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['CredentialDependencyEdge']} })
    edge_kind: str = Field(default=..., json_schema_extra = { "linkml_meta": {'domain_of': ['Requirement', 'CredentialDependencyEdge']} })
    derived_from_requirement: Optional[str] = Field(default=None, json_schema_extra = { "linkml_meta": {'domain_of': ['CredentialDependencyEdge']} })


# Model rebuild
# see https://pydantic-docs.helpmanual.io/usage/models/#rebuilding-a-model
BusinessEntity.model_rebuild()
Establishment.model_rebuild()
Process.model_rebuild()
Activity.model_rebuild()
ProductOffering.model_rebuild()
ChemicalHolding.model_rebuild()
EquipmentItem.model_rebuild()
AttributeDefinition.model_rebuild()
Fact.model_rebuild()
Scheme.model_rebuild()
SchemeVersion.model_rebuild()
Concept.model_rebuild()
Correspondence.model_rebuild()
ConceptMapping.model_rebuild()
ClassificationAssignment.model_rebuild()
CodeTranslation.model_rebuild()
RollupRule.model_rebuild()
Jurisdiction.model_rebuild()
JurisdictionProfile.model_rebuild()
Authority.model_rebuild()
LegalSource.model_rebuild()
Regime.model_rebuild()
Obligation.model_rebuild()
Determination.model_rebuild()
ProvenanceRecord.model_rebuild()
Credential.model_rebuild()
Requirement.model_rebuild()
CredentialDependencyEdge.model_rebuild()
